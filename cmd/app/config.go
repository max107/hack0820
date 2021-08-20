package main

import (
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type cfg struct {
	BotToken     string `env:"TELEGRAM_TOKEN"`
	AwsSecretKey string `env:"AWS_SECRET_ACCESS_KEY"`
	AwsAccessKey string `env:"AWS_ACCESS_KEY_ID"`
	AwsRegion    string `env:"AWS_REGION"`
	AwsBucket    string `env:"AWS_BUCKET"`
	JobURL       string `env:"JOB_URL"`
	StateURL     string `env:"STATE_URL"`
}

func LoadConfig(cfg interface{}, fileNames ...string) {
	if len(fileNames) == 0 {
		fileNames = []string{".env", ".env.local"}
	}

	var valid []string
	for _, f := range fileNames {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			continue
		}
		valid = append(valid, f)
	}
	_ = godotenv.Overload(valid...)
	_ = env.Parse(cfg)
	return
}
