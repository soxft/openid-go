package userutil

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
	"openid/config"
	"openid/library/tool"
	"openid/process/mysqlutil"
	"openid/process/redisutil"
	"regexp"
	"strings"
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

	userRedis := UserInfo{}
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
		UserId: userId,
		Iss:    config.C.Server.Title,
		Iat:    time.Now().Unix(),
		Jti:    Jti,
	})

	header := base64.StdEncoding.EncodeToString(headerJson)
	payload := base64.StdEncoding.EncodeToString(payloadJson)
	signature := header + "." + payload + "." + tool.Sha256(header+"."+payload, config.C.Jwt.Secret)
	return signature, nil
}

func getJti(user UserInfo) (string, error) {
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
	if _, err := _redis.Do("HMSET", redis.Args{}.Add(_redisKey).AddFlat(user)...); err != nil {
		log.Printf("[ERROR] getJti: %s", err.Error())
		return "", err
	}
	// 5分钟 登录后自动续期
	_, _ = _redis.Do("EXPIRE", _redisKey, 60*5)
	return _jti, nil
}

// GetJwtFromAuth
// 修改密码后 通过 payload 中的 IAT 和 JTI 来删除 redis 中的 JTI
func GetJwtFromAuth(Authorization string) string {
	reg, _ := regexp.Compile(`^Bearer\s+(.*)$`)
	if reg.MatchString(Authorization) {
		return reg.FindStringSubmatch(Authorization)[1]
	}
	return ""
}

func CheckJwt(jwt string) (UserInfo, error) {
	_jwt := strings.Split(jwt, ".")
	if len(_jwt) != 3 {
		return UserInfo{}, errors.New("jwt format error")
	}
	payloadJson, _ := base64.StdEncoding.DecodeString(_jwt[1])
	signature := _jwt[2]
	if tool.Sha256(_jwt[0]+"."+_jwt[1], config.C.Jwt.Secret) != signature {
		return UserInfo{}, errors.New("jwt signature error")
	}
	var payload JwtPayload
	err := json.Unmarshal(payloadJson, &payload)
	if err != nil {
		return UserInfo{}, errors.New("jwt payload error")
	}
	// check if JTI exists
	var userInfo UserInfo
	if userInfo, err = checkJti(payload.Jti, payload.Iat); err != nil {
		return UserInfo{}, err
	}

	return userInfo, nil
}

func checkJti(jti string, iat int64) (UserInfo, error) {
	_redis := redisutil.R.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)
	// get data from redis
	_redisKey := config.C.Redis.Prefix + ":jti:" + jti
	_redisData, err := redis.Values(_redis.Do("HGETALL", _redisKey))
	if err != nil {
		return UserInfo{}, errors.New("token expired")
	}
	if len(_redisData) == 0 {
		return UserInfo{}, errors.New("token expired")
	}
	// 解析数据
	var userInfo UserInfo
	if err = redis.ScanStruct(_redisData, &userInfo); err != nil {
		return UserInfo{}, errors.New("token expired")
	}
	// 判断是否有过期请求
	// 用户修改密码等操作后 会记录一个 xx:jti:expire:md5(username) 的 key 值为 修改密码时的时间戳, 用来与jwt中的iat进行比较
	expireTime, err := redis.Int64(_redis.Do("GET", config.C.Redis.Prefix+":jti:expire:"+tool.Md5(userInfo.Username)))
	if err != nil && err != redis.ErrNil {
		return UserInfo{}, errors.New("server error")
	}
	if iat < expireTime {
		_, _ = _redis.Do("DEL", _redisKey)
		return UserInfo{}, errors.New("token expired")
	}

	// 续期
	_, err = _redis.Do("EXPIRE", _redisKey, 60*60*24*14)
	return userInfo, nil
}
