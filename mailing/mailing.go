package mailing

import (
	"net/smtp"
)

func SendMail(recevier string) error {
	// Choose auth method and set it up
	from := "digitalonlinelibrary@gmail.com"
	pass := "bronetransporter"
	//to := "udartzev_sv@mail.ru"

	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{recevier}
	msg := []byte("To:" + recevier + "\r\n" +
		"Subject: Thanks for your interest in partnership!\r\n" +
		"\r\n" +
		"You are very cool author!\r\n")
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	if err != nil {
		return err
	}

	return nil
}
