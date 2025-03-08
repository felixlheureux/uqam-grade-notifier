package alert

import "net/smtp"

const smtp_server = "smtp.gmail.com"
const smtp_port = "587"

func SendEmail(pass, from, to, subject, body string) error {
	auth := smtp.PlainAuth("", from, pass, smtp_server)
	msg := []byte("From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body + "\n")

	return smtp.SendMail(smtp_server+":"+smtp_port, auth, from, []string{to}, msg)
}
