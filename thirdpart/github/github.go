package github

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/soxft/openid/config"
	"github.com/soxft/openid/library/toolutil"
	"log"
	"strconv"
)

const Platform = "Github"

type Github struct {
}

func New() *Github {
	return &Github{}
}

func (Github) Handler(c *gin.Context) (string, error) {
	// get access token
	resp, err := getAccessToken(c.Query("code"))
	if err != nil {
		log.Println("[ERROR] Github Login getAccessToken Error: ", err)
		return "", err
	}

	// get user id
	UserId, err := getUser(resp)
	if err != nil {
		log.Println("[ERROR] Github Login getUser Error: ", err)
		return "", err
	}
	return UserId, nil
}

func getAccessToken(code string) (string, error) {
	client := resty.New()

	accessToken := &AccessTokenStruct{}
	_, err := client.R().
		SetQueryParams(map[string]string{
			"client_id":     config.Github.ClientID,
			"client_secret": config.Github.ClientSecret,
			"code":          code,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&accessToken).
		Post("https://github.com/login/oauth/access_token")

	if accessToken.Error != "" {
		return "", fmt.Errorf("%s: %s", accessToken.Error, accessToken.ErrorDescription)
	} else if err != nil || accessToken.AccessToken == "" {
		return "", errors.New("request error, retry later")
	}
	return accessToken.AccessToken, nil
}

func getUser(accessToken string) (string, error) {
	client := resty.New()

	user := &UserStruct{}
	_, err := client.R().
		SetHeader("Authorization", fmt.Sprintf("token %s", accessToken)).
		SetHeader("Accept", "application/json").
		SetResult(&user).
		Get("https://api.github.com/user")

	if err != nil {
		return "", errors.New("request error, retry later")
	} else if user.ID == 0 {
		return "", errors.New(user.Message)
	}

	userId := toolutil.Sha256(strconv.Itoa(user.ID), "Github")
	return userId, nil
}
