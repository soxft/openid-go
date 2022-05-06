package userutil

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"log"
	"openid/config"
	"openid/library/tool"
	"openid/process/mysqlutil"
	"openid/process/redisutil"
	"time"
)

// GenerateJwt
// @description generate JWT token for user
func GenerateJwt(userId int) (string, error) {
	row, err := mysqlutil.D.Prepare("SELECT `id`,`username`,`email`,`lastTime`,`lastIp` FROM `account` WHERE `id` = ?")
	if err != nil {
		log.Printf("[ERROR] GenerateToken: %s", err.Error())
		return "", err
	}
	res := row.QueryRow(userId)
	if res.Err() == sql.ErrNoRows {
		return "", nil
	}

	userRedis := UserRedis{}
	_ = res.Scan(&userRedis.UserId, &userRedis.Username, &userRedis.Email, &userRedis.LastTime, &userRedis.LastIp)

	headerJson, _ := json.Marshal(JwtHeader{
		Alg: "HS256",
		Typ: "JWT",
	})

	var Jti string
	if Jti, err = getJti(userRedis); err != nil {
		log.Printf("[ERROR] GenerateToken: %s", err.Error())
		return "", err
	}
	payloadJson, _ := json.Marshal(JwtPayload{
		UserId:   userId,
		Username: userRedis.Username,
		Email:    userRedis.Email,
		Iss:      config.C.Server.Title,
		Iat:      time.Now().Unix(),
		Jti:      Jti,
	})

	header := base64.StdEncoding.EncodeToString(headerJson)
	payload := base64.StdEncoding.EncodeToString(payloadJson)
	signature := header + "." + payload + "." + tool.Sha256(header+"."+payload, config.C.Jwt.Secret)
	return signature, nil
}

func getJti(user UserRedis) (string, error) {
	JtiJson, _ := json.Marshal(map[string]string{
		"username": user.Username,
		"randStr":  tool.RandStr(32),
	})
	_redis := redisutil.R.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	_jti := tool.Md5(string(JtiJson))
	_redisKey := config.C.Redis.Prefix + ":jti:" + _jti
	_redisValue := map[string]interface{}{
		"userId":   user.UserId,
		"username": user.Username,
		"email":    user.Email,
		"lastTime": user.LastTime,
		"lastIp":   user.LastIp,
	}
	if _, err := _redis.Do("HMSET", redis.Args{}.Add(_redisKey).AddFlat(_redisValue)...); err != nil {
		log.Printf("[ERROR] getJti: %s", err.Error())
		return "", err
	}
	// 5分钟 登录后自动续期
	_, _ = _redis.Do("EXPIRE", _redisKey, 60*5)
	return _jti, nil
}

// 修改密码后 通过 payload 中的 IAT 和 JTI 来删除 redis 中的 JTI
