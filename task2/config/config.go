package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Configurations struct {
	DatabaseURL string `json:"database_url"`
}

func LoadConfigFromEnv() (*Configurations, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, nil
	}
	return &Configurations{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}, nil
}
