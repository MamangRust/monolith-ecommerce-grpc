package mailer

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Mailer struct {
	ctx      context.Context
	server   string
	port     int
	user     string
	password string
	tracer   trace.Tracer
}

func NewMailer(ctx context.Context, server string, port int, user string, password string) *Mailer {
	return &Mailer{
		ctx:      ctx,
		server:   server,
		port:     port,
		user:     user,
		password: password,
		tracer:   otel.Tracer("mailer"),
	}
}

func (m *Mailer) Send(to, subject, body string) error {
	_, span := m.tracer.Start(m.ctx, "SendEmail",
		trace.WithAttributes(
			attribute.String("email.recipient", to),
			attribute.String("email.subject", subject),
			attribute.String("smtp.server", m.server),
			attribute.Int("smtp.port", m.port),
		),
	)
	defer span.End()

	auth := smtp.PlainAuth("", m.user, m.password, m.server)

	headers := map[string]string{
		"From":         m.user,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": `text/html; charset="UTF-8"`,
	}

	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	addr := fmt.Sprintf("%s:%d", m.server, m.port)

	err := smtp.SendMail(addr, auth, m.user, []string{to}, msg.Bytes())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to send email")
	}

	return err
}
