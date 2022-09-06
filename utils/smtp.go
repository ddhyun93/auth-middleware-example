package utils

import (
	"fmt"
	"log"
	"net/smtp"
)

const (
	ID = "smtpinvestingdiary@gmail.com"
	PW = "fnnprrpizjgvjflk"
)

func SendAuthMail(email, code string) {
	auth := smtp.PlainAuth("", ID, PW, "smtp.gmail.com")
	from := ID
	to := []string{email}

	headerSubject := "Subject: investing-diary-email-auth-code\r\n"
	headerBlank := "\r\n"
	body := fmt.Sprintf("CODE : %s\r\n", code)
	msg := []byte(headerSubject + headerBlank + body)

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	log.Println(err)
}
