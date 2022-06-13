package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// create our Token structure
type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

/**
@brief: GenerateJWT will create a JWT (JSON web token) this token will save time role and mail from our user
*/
func GenerateJWT(email, role string) (string, error) {
	var mySigningKey = []byte("mysupersecretsecuritytoken")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Errorf("Something went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}
