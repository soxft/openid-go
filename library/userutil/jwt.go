package userutil

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/soxft/openid-go/app/model"
	"github.com/soxft/openid-go/config"
	"github.com/soxft/openid-go/library/toolutil"
	"github.com/soxft/openid-go/process/dbutil"
	"github.com/soxft/openid-go/process/redisutil"
	"gorm.io/gorm"
	"log"
	"regexp"
	"time"
)

// GetJwtFromAuth
// 从 Authorization 中获取JWT
func GetJwtFromAuth(Authorization string) string {
	reg, _ := regexp.Compile(`^Bearer\s+(.*)$`)
	if reg.MatchString(Authorization) {
		return reg.FindStringSubmatch(Authorization)[1]
	}
	return ""
}

// GenerateJwt
// @description generate JWT token for user
func GenerateJwt(userId int, clientIp string) (string, error) {
	var userInfo model.Account
	err := dbutil.D.Model(model.Account{}).Select("id, username, email, last_time, last_ip").Where(model.Account{ID: userId}).Take(&userInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	} else if err != nil {
		return "", err
	}

	timeNow := time.Now().Unix()

	uInfo := UserInfo{
		UserId:   userId,
		Username: userInfo.Username,
		Email:    userInfo.Email,
		LastTime: timeNow,
	}
	// update last login info
	setUserLastLogin(userInfo.ID, timeNow, clientIp)

	return jwt.NewWithClaims(jwt.SigningMethodHS512, JwtClaims{
		ID:       generateJti(uInfo),
		ExpireAt: time.Now().Add(24 * 30 * time.Hour).Unix(),
		IssuedAt: time.Now().Unix(),
		Issuer:   config.Server.Title,
		Username: userInfo.Username,
		UserId:   userId,
		Email:    userInfo.Email,
		LastTime: timeNow,
	}).SignedString([]byte(config.Jwt.Secret))
}

// CheckPermission
// @description check user permission
func CheckPermission(ctx context.Context, _jwt string) (UserInfo, error) {
	JwtClaims, err := JwtDecode(_jwt)
	if err != nil {
		return UserInfo{}, err
	}
	if checkJti(ctx, JwtClaims.ID) != nil {
		return UserInfo{}, ErrJwtExpired
	}
	return UserInfo{
		UserId:   JwtClaims.UserId,
		Username: JwtClaims.Username,
		Email:    JwtClaims.Email,
		LastTime: JwtClaims.LastTime,
	}, nil
}

// JwtDecode
// @description check JWT token
func JwtDecode(_jwt string) (JwtClaims, error) {
	token, err := jwt.ParseWithClaims(_jwt, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Jwt.Secret), nil
	})
	if err != nil {
		return JwtClaims{}, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return *claims, nil
	}

	return JwtClaims{}, errors.New("jwt token error")
}

// SetJwtExpire
// @description 标记JWT过期
func SetJwtExpire(c context.Context, _jwt string) error {
	JwtClaims, _ := JwtDecode(_jwt)
	_redis := redisutil.R

	ttl := JwtClaims.ExpireAt - time.Now().Unix()

	err := _redis.SetEx(c, getJwtExpiredKey(JwtClaims.ID), "1", time.Duration(ttl)).Err()
	if err != nil {
		log.Printf("[ERROR] SetJwtExpire: %s", err.Error())
		return errors.New("set jwt expire error")
	}

	return nil
}

// checkJti
func checkJti(ctx context.Context, jti string) error {
	_redis := redisutil.R

	expired, err := _redis.Exists(ctx, getJwtExpiredKey(jti)).Result()
	if err != nil {
		log.Printf("[ERROR] checkJti: %s", err.Error())
		return errors.New("check jti error")
	}
	if expired == 1 {
		return ErrJwtExpired
	}

	return nil
}

// generateJti 创建 Jti
func generateJti(user UserInfo) string {
	JtiJson, _ := json.Marshal(map[string]string{
		"username": user.Username,
		"randStr":  toolutil.RandStr(32),
		"time":     time.Now().Format("150405"),
	})
	_jti := toolutil.Md5(string(JtiJson))
	return _jti
}

func setUserLastLogin(userId int, lastTime int64, lastIp string) {
	dbutil.D.Model(&model.Account{}).Where(model.Account{ID: userId}).Updates(&model.Account{LastTime: lastTime, LastIp: lastIp})
}

// getJwtExpiredKey
func getJwtExpiredKey(jti string) string {
	return config.RedisPrefix + ":jti:expired:" + jti
}
