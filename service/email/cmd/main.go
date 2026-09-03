package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MamangRust/monolith-ecommerce-grpc-email/internal/config"
	"github.com/MamangRust/monolith-ecommerce-grpc-email/internal/handler"
	"github.com/MamangRust/monolith-ecommerce-grpc-email/internal/mailer"
	"github.com/MamangRust/monolith-ecommerce-grpc-email/internal/metrics"
	"github.com/MamangRust/monolith-ecommerce-pkg/database"
	"github.com/MamangRust/monolith-ecommerce-pkg/dotenv"
	"github.com/MamangRust/monolith-ecommerce-pkg/emailretry"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	otel_pkg "github.com/MamangRust/monolith-ecommerce-pkg/otel"
	"github.com/MamangRust/monolith-ecommerce-pkg/outbox"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	if err := dotenv.Viper(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	otelCfg := otel_pkg.Config{
		ServiceName:          "email-service",
		ServiceVersion:       "1.0.0",
		Environment:          viper.GetString("ENV"),
		Endpoint:             viper.GetString("OTEL_EXPORTER_OTLP_ENDPOINT"),
		Insecure:             true,
		EnableRuntimeMetrics: true,
	}

	if otelCfg.Endpoint == "" {
		otelCfg.Endpoint = "otel-collector:4317"
	}

	telemetry := otel_pkg.NewTelemetry(otelCfg)

	if err := telemetry.Init(ctx); err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}

	defer func() {
		if err := telemetry.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down telemetry: %v", err)
		}
	}()

	logger, err := logger.NewLogger("email-service", telemetry.GetLogger())
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	cfg := config.Config{
		KafkaBrokers: []string{viper.GetString("KAFKA_BROKERS")},
		SMTPServer:   viper.GetString("SMTP_SERVER"),
		SMTPPort:     viper.GetInt("SMTP_PORT"),
		SMTPUser:     viper.GetString("SMTP_USER"),
		SMTPPass:     viper.GetString("SMTP_PASS"),
		MaxRetries:   viper.GetInt("EMAIL_MAX_RETRIES"),
		RetryBackoff: viper.GetDuration("EMAIL_RETRY_BACKOFF"),
	}
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = emailretry.DefaultMaxAttempts
	}
	if cfg.RetryBackoff <= 0 {
		cfg.RetryBackoff = emailretry.DefaultBackoff
	}

	// Register OTel metric instruments after the SDK is initialized so they
	// are bound to the real meter provider (exported via OTLP), not the noop
	// default from package init.
	metrics.Register()

	m := mailer.NewMailer(
		ctx,
		cfg.SMTPServer,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPass,
	)

	// The consumer inbox lives in PostgreSQL, so the email service now requires
	// a database connection at startup. If the database is unreachable the
	// service refuses to start rather than silently losing the idempotency
	// guarantee.
	dbPool, err := database.NewClient(logger)
	if err != nil {
		logger.Fatal("Failed to connect to database for consumer inbox", zap.Error(err))
	}
	defer dbPool.Close()

	inbox, err := outbox.NewPostgresInbox(dbPool)
	if err != nil {
		logger.Fatal("Failed to initialize consumer inbox", zap.Error(err))
	}

	myKafka := kafka.NewKafka(logger, cfg.KafkaBrokers)

	h := handler.NewEmailHandlerWithInbox(ctx, logger, m, inbox, "email-service-group", myKafka, cfg.RetryBackoff)

	err = myKafka.StartConsumersWithContextManualCommit(ctx, []string{
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

	retryH := handler.NewRetryHandler(ctx, logger, m, inbox, "email-service-group", myKafka, cfg.MaxRetries, cfg.RetryBackoff)
	if err := myKafka.StartConsumersWithContextManualCommit(ctx, []string{emailretry.RetryTopic}, emailretry.RetryGroup, retryH); err != nil {
		log.Fatalf("Error starting retry consumer: %v", err)
	}

	logger.Info("Email service started", zap.String("retry_topic", emailretry.RetryTopic), zap.String("dlq_topic", emailretry.DLQTopic))

	<-ctx.Done()
	logger.Info("Shutting down email service")
	if err := myKafka.Close(); err != nil {
		logger.Error("Failed to close Kafka resources", zap.Error(err))
	}
}
