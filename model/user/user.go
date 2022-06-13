package user

import (
	"../../database"
	"crypto/rand"
	"fmt"
	"github.com/jinzhu/gorm"
)

// create a struct named User
type User struct {
	gorm.Model
	Pseudo        string `gorm:"column:pseudo" json:"pseudo"`
	Email         string `gorm:"column:email;unique" json:"email"`
	Password      string `gorm:"column:password;type:varchar(255)" json:"-"`
	TokenEmail    string `gorm:"column:token_email" json:"-"`
	EmailValidate bool   `gorm:"column:email_validate" json:"email_validate"`
	Role          string `gorm:"column:role" json:"role"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Register struct {
	Pseudo   string `json:"pseudo"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

/**
@brief: TokenGenerator generate a token(used in validation mail)
*/
func TokenGenerator(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

/**
@brief: NewUser will send our user and his parameters in database
*/
func NewUser(user *User) (User, error) {
	db := database.DBConn
	newUser := User{Email: user.Email, Pseudo: user.Pseudo, Password: user.Password, TokenEmail: user.TokenEmail, EmailValidate: user.EmailValidate, Role: user.Role}
	err := db.Create(&newUser)
	return newUser, err.Error
}

/**
@brief: GetUserByToken identify our user from his token
*/
func GetUserByToken(token string) (User, bool) {
	db := database.DBConn
	user := User{}
	db.Where("token_email = ? and email_validate = ?", token, false).Find(&user)
	if user.Email == "" {
		return user, false
	}
	return user, true
}

/**
@brief: GetUserByEmail search user in our database from his email address
*/
func GetUserByEmail(email string) (User, bool) {
	db := database.DBConn
	user := User{}
	db.Where("email = ?", email).Find(&user)
	if user.Email == "" {
		return user, false
	}
	return user, true
}

/**
@brief: GetUserById search user in our database from his ID
*/
func GetUserById(id int) User {
	db := database.DBConn
	user := User{}
	db.Where("id = ?", id).Find(&user)
	return user
}

/**
@brief: updateUser will update user parameter
*/
func UpdateUser(user User, updateUser User) User {
	db := database.DBConn
	db.Model(&user).Update(updateUser)
	return user
}

func DeleteUser(user User) {
	db := database.DBConn
	db.Delete(user)
}
