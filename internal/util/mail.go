package util

import (
	"bytes"
	htmlTemplate "html/template"
	textTemplate "text/template"

	"github.com/Adedunmol/mycart/internal/config"
	"github.com/Adedunmol/mycart/internal/logger"
	"gopkg.in/gomail.v2"
)

// type Template string

// const (
// 	TemplateVerification   Template = "verification"
// 	TemplateForgotPassword Template = "forgot-password"
// 	TemplatePurchase       Template = "purchase"
// )

func SendMail(to string, subject string, html string, plain string) {
	m := gomail.NewMessage()

	m.SetHeader("From", config.EnvConfig.EmailUsername)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", html)
	m.SetBody("text/plain", plain)

	d := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, config.EnvConfig.EmailUsername, config.EnvConfig.EmailPassword)

	if err := d.DialAndSend(m); err != nil {
		logger.Error.Println("could not send mail")
	}
}

func SendMailWithTemplate(templateFile string, to string, subject string, locals struct{}) {
	m := gomail.NewMessage()

	m.SetHeader("From", config.EnvConfig.EmailUsername)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	var html *htmlTemplate.Template
	var text *textTemplate.Template
	var err error

	switch templateFile {
	case "verification":
		html, err = html.ParseFiles("../email-templates/verification/verification.html")
		if err != nil {
			logger.Error.Println("could not parse html file")
			return
		}

		text, err = text.ParseFiles("../email-templates/verification/verification.txt")
		if err != nil {
			logger.Error.Println("could not parse text file")
			return
		}
	}

	htmlBuff := new(bytes.Buffer)
	html.Execute(htmlBuff, locals)

	textBuff := new(bytes.Buffer)
	text.Execute(htmlBuff, locals)

	m.SetBody("text/html", htmlBuff.String())
	m.SetBody("text/plain", textBuff.String())

	d := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, config.EnvConfig.EmailUsername, config.EnvConfig.EmailPassword)

	if err := d.DialAndSend(m); err != nil {
		logger.Error.Println("could not send mail")
	}
}
