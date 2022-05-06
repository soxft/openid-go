package tool

import "regexp"

func IsEmail(email string) bool {
	b, _ := regexp.MatchString("^([a-z0-9_.-]+)@([da-z.-]+).([a-z.]{2,6})$", email)
	return b
}

func IsUserName(id string) bool {
	b, _ := regexp.MatchString("^[0-9a-zA-Z]{5,30}$", id)
	return b
}

func IsPassword(password string) bool {
	if len(password) < 6 || len(password) > 128 {
		return false
	}
	return true
}
