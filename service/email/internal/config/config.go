package config

type Config struct {
	KafkaBrokers []string
	SMTPServer   string
	SMTPPort     int
	SMTPUser     string
	SMTPPass     string
}
