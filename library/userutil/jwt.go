package userutil

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/soxft/openid/config"
	"github.com/soxft/openid/library/toolutil"
	"github.com/soxft/openid/process/dbutil"
	"github.com/soxft/openid/process/redisutil"
	"gorm.io/gorm"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GenerateJwt
// @description generate JWT token for user
func GenerateJwt(userId int, clientIp string) (string, error) {
	var userInfo dbutil.Account
	err := dbutil.D.Model(dbutil.Account{}).Select("id, username, email, last_time, last_ip").Where(dbutil.Account{ID: userId}).Take(&userInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	userRedis := UserInfo{
		UserId:   userId,
		Username: userInfo.Username,
		Email:    userInfo.Email,
	}
	userLast := UserLastInfo{
		LastIp:   userInfo.LastIp,
		LastTime: userInfo.LastTime,
	}
	_ = setUserBaseInfo(userRedis.UserId, userLast)

	// update last login info
	dbutil.D.Model(&dbutil.Account{}).Where(dbutil.Account{ID: userId}).Updates(&dbutil.Account{LastTime: time.Now().Unix(), LastIp: clientIp})

	headerJson, _ := json.Marshal(JwtHeader{
		Alg: "HS256",
		Typ: "JWT",
	})

	var Jti string
	if Jti, err = generateJti(userRedis); err != nil {
		log.Printf("[ERROR] GenerateToken: %s", err.Error())
		return "", err
	}
	payloadJson, _ := json.Marshal(JwtPayload{
		UserId: userId,
		Iss:    config.Server.Title,
		Iat:    time.Now().Unix(),
		Jti:    Jti,
	})

	header := base64.StdEncoding.EncodeToString(headerJson)
	payload := base64.StdEncoding.EncodeToString(payloadJson)
	signature := header + "." + payload + "." + toolutil.Sha256(header+"."+payload, config.Jwt.Secret)
	return signature, nil
}

func generateJti(user UserInfo) (string, error) {
	JtiJson, _ := json.Marshal(map[string]string{
		"username": user.Username,
		"randStr":  toolutil.RandStr(32),
	})
	_redis := redisutil.R.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	_jti := toolutil.Md5(string(JtiJson))
	_redisKey := config.RedisPrefix + ":jti:" + _jti
	if _, err := _redis.Do("HMSET", redis.Args{}.Add(_redisKey).AddFlat(user)...); err != nil {
		log.Printf("[ERROR] getJti: %s", err.Error())
		return "", err
	}
	// 默认5分钟 登录后自动续期
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

// CheckJwt
// @description check JWT token
// @param 用户每次请求是验证JWT token
func CheckJwt(jwt string) (UserInfo, error) {
	_jwt := strings.Split(jwt, ".")
	if len(_jwt) != 3 {
		return UserInfo{}, errors.New("jwt format error")
	}
	payloadJson, _ := base64.StdEncoding.DecodeString(_jwt[1])
	signature := _jwt[2]
	if toolutil.Sha256(_jwt[0]+"."+_jwt[1], config.Jwt.Secret) != signature {
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
	// 续期
	go setExpire(":jti:"+payload.Jti, ":user:last:"+strconv.Itoa(userInfo.UserId))

	return userInfo, nil
}

// DelJti
// @description check if JTI exists
func DelJti(jwt string) error {
	_redis := redisutil.R.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)
	// get jti from jwt
	_jwt := strings.Split(jwt, ".")
	if len(_jwt) != 3 {
		return errors.New("jwt format error")
	}

	if payloadJson, err := base64.StdEncoding.DecodeString(_jwt[1]); err != nil {
		return errors.New("jwt payload error")
	} else {
		var payload JwtPayload
		if err = json.Unmarshal(payloadJson, &payload); err != nil {
			return errors.New("jwt payload error")
		}
		_, err = _redis.Do("DEL", config.RedisPrefix+":jti:"+payload.Jti)
		return err
	}
}

// SetUserJwtExpire
// @description 修改密码等操作后 使用户所有的jwt token过期
func SetUserJwtExpire(username string, expire int64) error {
	_redis := redisutil.R.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)

	_redisKey := config.RedisPrefix + ":jti:expire:" + toolutil.Md5(username)
	_, err := _redis.Do("SET", _redisKey, expire)
	if err != nil {
		log.Printf("[ERROR] SetUserJwtExpire: %s", err.Error())
		return errors.New("set user jwt expire error")
	}
	return nil
}

// checkJti
func checkJti(jti string, iat int64) (UserInfo, error) {
	_redis := redisutil.R.Get()
	defer func(_redis redis.Conn) {
		_ = _redis.Close()
	}(_redis)
	// get data from redis
	_redisKey := config.RedisPrefix + ":jti:" + jti
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
	// 判断是否有过期请求 TODO
	// 用户修改密码等操作后 会记录一个 xx:jti:expire:md5(username) 的 key 值为 修改密码时的时间戳, 用来与jwt中的iat进行比较
	expireTime, err := redis.Int64(_redis.Do("GET", config.RedisPrefix+":jti:expire:"+toolutil.Md5(userInfo.Username)))
	if err != nil && err != redis.ErrNil {
		return UserInfo{}, errors.New("server error")
	}
	if iat < expireTime {
		_, _ = _redis.Do("DEL", _redisKey)
		return UserInfo{}, errors.New("token expired")
	}

	return userInfo, nil
}

func setUserBaseInfo(userId int, user UserLastInfo) error {
	_redis := redisutil.R.Get()
	_redisKey := config.RedisPrefix + ":user:last:" + strconv.Itoa(userId)
	_, err := _redis.Do("HMSET", redis.Args{}.Add(_redisKey).AddFlat(user)...)
	_ = _redis.Close()
	return err
}

func setExpire(redisKey ...string) {
	_redis := redisutil.R.Get()
	for _, key := range redisKey {
		_redisKey := config.RedisPrefix + key
		_, _ = _redis.Do("EXPIRE", _redisKey, 60*60*24*64)
	}
	_ = _redis.Close()
}
