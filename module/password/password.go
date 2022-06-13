package password

import (
	"golang.org/x/crypto/bcrypt"
)

/**
@brief: hash our password ( we can't allow to have a "clear" password wrote in the api)
*/
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

/**
@brief: CheckPasswordHash will check if our password and our hash match
*/
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
