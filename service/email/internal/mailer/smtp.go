package mailer

import (
	"bytes"
	"fmt"
	"net/smtp"
)

type Mailer struct {
	Server   string
	Port     int
	User     string
	Password string
}

func (m *Mailer) Send(to string, subject string, body string) error {
	auth := smtp.PlainAuth("", m.User, m.Password, m.Server)

	headers := make(map[string]string)
	headers["From"] = m.User
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	addr := fmt.Sprintf("%s:%d", m.Server, m.Port)
	return smtp.SendMail(addr, auth, m.User, []string{to}, msg.Bytes())
}
