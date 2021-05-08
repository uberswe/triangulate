package triangulate

import (
	"fmt"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
)

func sendEmail(receiver string, subject string, content string) {
	e := email.NewEmail()
	e.From = fmt.Sprintf("Triangulate.xyz <%s>", fromEmail)
	e.To = []string{receiver}
	e.Subject = subject
	e.Text = []byte(content)
	err := e.Send(fmt.Sprintf("%s:%s", smtpHost, smtpPort), smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Email Sent Successfully!")
}
