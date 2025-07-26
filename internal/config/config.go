package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Server Server
	Enode  Enode
	Redis  Redis
}

type Server struct {
	Port string `env:"PORT,required"`
}

type Enode struct {
	ClientID     string `env:"ENODE_CLIENT_ID,required"`
	ClientSecret string `env:"ENODE_CLIENT_SECRET,required"`
	OAuthBaseURL string `env:"ENODE_OAUTH_URL,required"`
	ApiURL       string `env:"ENODE_API_URL,required"`
}

type Redis struct {
	Host     string `env:"REDIS_HOST,required"`
	Port     string `env:"REDIS_PORT,required"`
	Password string `env:"REDIS_PASSWORD,required"`
	DB       int    `env:"REDIS_DB,required"`
}

func LoadConfig(envFile string) (*Config, error) {
	var cfg Config
	_ = godotenv.Load(envFile)

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse env: %v", err)
	}

	return &cfg, nil
}
