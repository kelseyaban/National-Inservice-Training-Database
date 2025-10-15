// Filename: internal/mailer/mailer.go
package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)

// don't need a separate server for serving static files
//
//go:embed templates/*
var templateFS embed.FS // embed the files from templates into our program

type Mailer struct {
	dialer *mail.Dialer // connection to the SMTP server
	sender string       // who is sending the email
}

// Configure a SMTP connection instance using our credentials from Mailtrap
func New(host string, port int, username, password, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

// Send the email to the user. The data parameter is for the dynamic data
// to inject into the template
func (m Mailer) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)

	if err != nil {
		return err
	}
	println("Template parsed successfully", templateFile)

	// fill in the subject part
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	// fill in the plainBody part
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	// fill in the htmlBody part
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	// Craft the message from the parts above
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	// send the message. Try to do this at most 3 times before giving up
	for i := 1; i <= 3; i++ {
		err = m.dialer.DialAndSend(msg)
		// If everything worked, return nil.
		if err == nil {
			return nil
		}

		// If it didn't work, sleep for a short time and retry.
		time.Sleep(500 * time.Millisecond)
	}

	return err // give up
}
