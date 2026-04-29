package dotenv

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func Viper() error {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	var configFile string
	useConfigFile := true

	switch env {
	case "docker":
		configFile = "/app/docker.env"
	case "production":
		configFile = "/app/production.env"
	case "kubernetes", "test":
		useConfigFile = false
	default:
		configFile = ".env"
	}

	viper.AutomaticEnv()

	if useConfigFile {
		viper.SetConfigFile(configFile)

		err := viper.ReadInConfig()
		if err != nil {
			return fmt.Errorf("error reading config file %s: %w", configFile, err)
		}
	}

	return nil
}
