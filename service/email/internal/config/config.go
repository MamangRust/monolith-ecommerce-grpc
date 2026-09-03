package config

import "time"

type Config struct {
	KafkaBrokers []string
	SMTPServer   string
	SMTPPort     int
	SMTPUser     string
	SMTPPass     string
	// MaxRetries is the number of SMTP attempts (including the first) before
	// an event is dead-lettered (Phase 4).
	MaxRetries int
	// RetryBackoff is the base delay before the next attempt; the actual delay
	// grows linearly with the attempt number (emailretry.Backoff).
	RetryBackoff time.Duration
}
