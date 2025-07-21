package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db   DbConfig
	LLM  LLMConfig
	Auth AuthConfig
}

type DbConfig struct {
	Dsn string
}

type LLMConfig struct {
	URL    string
	ApiKey string
}

type AuthConfig struct {
	Secret string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default config")
	}
	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},

		LLM: LLMConfig{
			URL:    os.Getenv("URL"),
			ApiKey: os.Getenv("ApiKey"),
		},

		Auth: AuthConfig{
			Secret: os.Getenv("TOKEN"),
		},
	}

}
