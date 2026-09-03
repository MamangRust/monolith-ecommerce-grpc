// Package event provides the standard event envelope shared by producers and
// consumers of email events (Phase 2 — Standard Event Envelope).
package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// SchemaVersion is the current envelope schema version. Consumers reject
// events whose schema_version differs.
const SchemaVersion = 1

// EmailEnvelope is the standard envelope for email events. Every producer of
// an email topic must publish this shape so consumers can validate the event
// before sending anything to SMTP.
type EmailEnvelope struct {
	EventID       string         `json:"event_id"`
	SchemaVersion int            `json:"schema_version"`
	EventType     string         `json:"event_type"`
	OccurredAt    string         `json:"occurred_at"`
	Email         string         `json:"email"`
	Subject       string         `json:"subject"`
	Body          string         `json:"body"`
	Metadata      map[string]any `json:"metadata,omitempty"`
}

// NewEmail builds an envelope with a fresh UUID event_id and the current time.
func NewEmail(eventType, email, subject, body string) EmailEnvelope {
	return EmailEnvelope{
		EventID:       uuid.NewString(),
		SchemaVersion: SchemaVersion,
		EventType:     eventType,
		OccurredAt:    time.Now().UTC().Format(time.RFC3339),
		Email:         email,
		Subject:       subject,
		Body:          body,
	}
}

// MarshalEmail returns the JSON bytes of an enveloped email event.
func MarshalEmail(eventType, email, subject, body string) ([]byte, error) {
	return json.Marshal(NewEmail(eventType, email, subject, body))
}

// IsValid reports whether the envelope satisfies the consumer contract:
// a non-empty event_id, schema_version == 1, a non-empty event_type and the
// email/subject/body fields. Invalid events must never reach SMTP.
func (e EmailEnvelope) IsValid() bool {
	return e.EventID != "" &&
		e.SchemaVersion == SchemaVersion &&
		e.EventType != "" &&
		e.Email != "" &&
		e.Subject != "" &&
		e.Body != ""
}
