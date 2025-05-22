package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	EmailSent = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "email_sent_total",
		Help: "Total emails sent successfully",
	})

	EmailFailed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "email_failed_total",
		Help: "Total emails failed",
	})
)

func Register() {
	prometheus.MustRegister(EmailSent, EmailFailed)
}
