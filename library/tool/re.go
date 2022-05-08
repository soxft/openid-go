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
	if len(password) < 8 || len(password) > 64 {
		return false
	}
	return true
}

func IsDomain(domain string) bool {
	b, _ := regexp.MatchString("^[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]\\.[a-zA-Z.]{2,6}$", domain)
	return b
}
