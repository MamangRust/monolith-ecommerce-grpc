package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/MamangRust/monolith-ecommerce-grpc-email/internal/config"
	"github.com/MamangRust/monolith-ecommerce-grpc-email/internal/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-email/internal/mailer"
	"github.com/MamangRust/monolith-ecommerce-grpc-email/internal/metrics"
	"github.com/MamangRust/monolith-ecommerce-pkg/dotenv"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	logger, err := logger.NewLogger("email-service")
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	if err := dotenv.Viper(); err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
	}
	ctx := context.Background()

	cfg := config.Config{
		KafkaBrokers: []string{viper.GetString("KAFKA_BROKERS")},
		SMTPServer:   viper.GetString("SMTP_SERVER"),
		SMTPPort:     viper.GetInt("SMTP_PORT"),
		SMTPUser:     viper.GetString("SMTP_USER"),
		SMTPPass:     viper.GetString("SMTP_PASS"),
	}

	metricsAddr := fmt.Sprintf(":%s", viper.GetString("METRIC_EMAIL_ADDR"))

	metrics.Register()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(metricsAddr, nil))
	}()

	m := mailer.NewMailer(
		ctx,
		cfg.SMTPServer,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPass,
	)

	h := handler.NewEmailHandler(ctx, logger, m)

	myKafka := kafka.NewKafka(logger, cfg.KafkaBrokers)

	err = myKafka.StartConsumers([]string{
		"email-service-topic-auth-register",
		"email-service-topic-auth-forgot-password",
		"email-service-topic-auth-verify-code-success",
		"email-service-topic-merchant-create",
		"email-service-topic-merchant-update-status",
		"email-service-topic-merchant-document-create",
		"email-service-topic-merchant-document-update-status",
		"email-service-topic-transaction-create",
	}, "email-service-group", h)

	if err != nil {
		log.Fatalf("Error starting consumer: %v", err)
	}
	select {}
}
