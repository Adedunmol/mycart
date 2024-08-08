package util

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	htmlTemplate "html/template"
	"os"
	"path/filepath"
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

	m.SetHeader("From", config.EnvConfig.EmailSender)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", html)
	m.SetBody("text/plain", plain)

	d := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, config.EnvConfig.EmailUsername, config.EnvConfig.EmailPassword)

	if err := d.DialAndSend(m); err != nil {
		logger.Logger.Error("could not send mail")
	}
}

func SendMailWithTemplate(templateFile string, to string, subject string, locals interface{}, attachment string) {
	m := gomail.NewMessage()

	m.SetHeader("From", config.EnvConfig.EmailSender)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	var html *htmlTemplate.Template
	var text *textTemplate.Template
	// var err error

	currentDirectory, _ := os.Getwd()

	switch templateFile {
	case "verification":
		pathToHtmlFile := filepath.Join(currentDirectory, "internal", "email-templates", "verification", "verification.html")
		if _, err := os.Stat(pathToHtmlFile); err == nil {
			// path/to/whatever exists
			fmt.Println("html file does exist")
			html = htmlTemplate.Must(htmlTemplate.ParseFiles(pathToHtmlFile))
		} else if errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does *not* exist
			logger.Logger.Error("html file does not exist")
			return
		}

		pathToTextFile := filepath.Join(currentDirectory, "internal", "email-templates", "verification", "verification.txt")
		if _, err := os.Stat(pathToTextFile); err == nil {
			// path/to/whatever exists
			fmt.Println("text file does exist")
			text = textTemplate.Must(textTemplate.ParseFiles(pathToTextFile))
		} else if errors.Is(err, os.ErrNotExist) {
			logger.Logger.Error("text file does not exist")
			return
		}

	case "purchase":
		pathToHtmlFile := filepath.Join(currentDirectory, "internal", "email-templates", "purchase", "purchase.html")
		if _, err := os.Stat(pathToHtmlFile); err == nil {
			// path/to/whatever exists
			fmt.Println("html file does exist")
			html = htmlTemplate.Must(htmlTemplate.ParseFiles(pathToHtmlFile))
		} else if errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does *not* exist
			logger.Logger.Error("html file does not exist")
			return
		}

		pathToTextFile := filepath.Join(currentDirectory, "internal", "email-templates", "purchase", "purchase.txt")
		if _, err := os.Stat(pathToTextFile); err == nil {
			// path/to/whatever exists
			fmt.Println("text file does exist")
			html = htmlTemplate.Must(htmlTemplate.ParseFiles(pathToHtmlFile))
		} else if errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does *not* exist
			logger.Logger.Error("text file does not exist")
			return
		}
	}
	textBuff := new(bytes.Buffer)
	text.Execute(textBuff, locals)

	m.AddAlternative("text/plain", textBuff.String())

	htmlBuff := new(bytes.Buffer)
	html.Execute(htmlBuff, locals)

	m.SetBody("text/html", htmlBuff.String())

	if attachment != "" {
		m.Attach(attachment)
	}

	d := gomail.NewDialer("sandbox.smtp.mailtrap.io", 2525, config.EnvConfig.EmailUsername, config.EnvConfig.EmailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		logger.Logger.Error(err.Error())
		logger.Logger.Error("could not send mail")
	}

	logger.Logger.Info("email sent to " + to)
}
