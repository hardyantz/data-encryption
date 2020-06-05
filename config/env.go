package config

import (
	"fmt"
	"os"

	config "github.com/joho/godotenv"
)

func Load(envPath string) error {
	err := config.Load(envPath)
	if err != nil {
		return fmt.Errorf(".env is not loaded properly")
	}

	if _, ok := os.LookupEnv("DB_HOST"); !ok {
		return fmt.Errorf("DB host is not loaded")
	}

	if _, ok := os.LookupEnv("DB_NAME"); !ok {
		return fmt.Errorf("DB name is not loaded")
	}

	if _, ok := os.LookupEnv("DB_USER"); !ok {
		return fmt.Errorf("DB user is not loaded")
	}

	if _, ok := os.LookupEnv("DB_PASS"); !ok {
		return fmt.Errorf("DB pass is not loaded")
	}

	return nil
}
