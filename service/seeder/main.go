package main

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-pkg/database"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/database/seeder"
	"github.com/MamangRust/monolith-ecommerce-pkg/dotenv"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	logger, err := logger.NewLogger("seeder")
	if err != nil {
		logger.Fatal("Failed to initialize logger", zap.Error(err))
	}

	if err := dotenv.Viper(); err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
	}

	conn, err := database.NewClient(logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	DB := db.New(conn)

	ctx := context.Background()

	hash := hash.NewHashingPassword()

	db_seeder := viper.GetString("DB_SEEDER")

	if db_seeder == "true" {
		logger.Debug("Seeding database", zap.String("seeder", db_seeder))
		seeder := seeder.NewSeeder(seeder.Deps{
			Db:     DB,
			Hash:   hash,
			Ctx:    ctx,
			Logger: logger,
		})
		if err := seeder.Run(); err != nil {
			logger.Fatal("Failed to run seeder", zap.Error(err))
		}
	}

	logger.Info("Database seeded successfully")
}
