package util

import (
	"github.com/Adedunmol/mycart/internal/config"
	"github.com/Adedunmol/mycart/internal/logger"
	"gopkg.in/gomail.v2"
)

func SendMail(to string, subject string, html string, plain string) {
	m := gomail.NewMessage()

	m.SetHeader("From", config.EnvConfig.EmailUsername)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", html)
	m.SetBody("text/plain", plain)

	d := gomail.NewDialer("smtp.gmail.com", 587, config.EnvConfig.EmailUsername, config.EnvConfig.EmailPassword)

	if err := d.DialAndSend(m); err != nil {
		logger.Error.Println("could not send mail")
	}
}
