package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

var (
	from     = "golang.sendemail.projectforum@gmail.com"
	password = "10012002Lo"
	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
)

/**
@brief: SendMail will send mail to user
*/
func SendMail(to []string, data interface{}, tTemplate string) {

	auth := smtp.PlainAuth("", from, password, smtpHost)
	t, _ := template.ParseFiles("public/emailTemplate/" + tTemplate)

	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))
	t.Execute(&body, data)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Email Sent!")
	return
}
