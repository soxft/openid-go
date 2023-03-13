package userutil

import "golang.org/x/crypto/bcrypt"

// GeneratePwd hash password
func GeneratePwd(password string) (string, error) {
	str, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(str), err
}

// CheckPwd check if password is correct
func CheckPwd(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
