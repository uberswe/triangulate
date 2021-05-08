package triangulate

import (
	"fmt"
	"net/smtp"
)

func sendEmail(receiver string, subject string, content string) {

	// Receiver email address.
	to := []string{
		receiver,
	}

	// Message.
	message := []byte(fmt.Sprintf(`Subject: %s

		%s`, subject, content))

	// Authentication.
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}
