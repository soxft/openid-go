package thirdpart

import (
	"errors"
	"github.com/soxft/openid/app/model"
	"github.com/soxft/openid/library/userutil"
	"github.com/soxft/openid/process/dbutil"
	"gorm.io/gorm"
	"log"
)

// Handler
// 登录处理
// @param identify 用户标识
// @param platform 平台
// @param ip 用户IP
// @return string jwt_token
func Handler(identify string, platform, ip string) (string, error) {
	// get user_id by identify
	userId, err := getUserIdByIdentify(identify, platform)
	if err != nil {
		return "", err
	}
	var jwt string
	if jwt, err = userutil.GenerateJwt(userId, ip); err != nil {
		return "", errors.New("system error, retry later")
	}
	return jwt, nil
}

func getUserIdByIdentify(identify string, platform string) (int, error) {
	// get user_id by identify
	var userId int
	err := dbutil.D.Model(model.Third{}).Select("user_id").
		Where(model.Third{
			Identify: identify,
			Platform: platform,
		}).
		Take(&userId).Error

	if err == gorm.ErrRecordNotFound {
		// 别问为什么要这样写, 问就是懒
		return 0, errors.New("(懒) 请先登录后绑定该平台")
	} else if err != nil {
		log.Println("[ERROR] getUserIdByIdentify error:", err)
		return 0, errors.New("system error, retry later")
	}
	return userId, nil
}
