package handler

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-ecommerce-grpc-email/internal/mailer"
	"github.com/MamangRust/monolith-ecommerce-grpc-email/internal/metrics"
)

type EmailHandler struct {
	Mailer *mailer.Mailer
}

func (h *EmailHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *EmailHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *EmailHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var payload map[string]interface{}
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		email := payload["email"].(string)
		subject := payload["subject"].(string)
		body := payload["body"].(string)

		err := h.Mailer.Send(email, subject, body)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
			metrics.EmailFailed.Inc()
		} else {
			metrics.EmailSent.Inc()
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
