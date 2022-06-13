package controller

import (
	"../module/email"
)

// work in progress
func SendEmail() {
	email.SendMail([]string{"louissasse0@gmail.com"}, struct {
		UrlServer string
		Token     string
	}{UrlServer: "https://localhost:8080", Token: "dlkfsdkjfjlsnfksdfdklqdklqskld"}, "verifyEmail.html")
}
