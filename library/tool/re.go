package tool

import "regexp"

func IsEmail(email string) bool {
	b, _ := regexp.MatchString("^([a-z0-9_.-]+)@([da-z.-]+).([a-z.]{2,6})$", email)
	return b
}

func IsID(id string) bool {
	b, _ := regexp.MatchString("^[0-9a-zA-Z]{5,30}$", id)
	return b
}
