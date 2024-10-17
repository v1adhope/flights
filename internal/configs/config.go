package configs

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Postgres Postgres
		Srv      Srv
	}

	Postgres struct {
		ConnStr string `env-required:"true" env:"SERVICE_POSTGRES_CONN_STR"`
	}

	Srv struct {
		Mode            string        `env-required:"true" env:"SERVICE_SRV_MODE"`
		ShutdownTimeout time.Duration `env-required:"true" env:"SERVICE_SRV_SHUTDOWN_TIMEOUT"`
	}
)

var Global Config

func MustConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("config: can't load envs from .env: %v", err)
	}

	if err := cleanenv.ReadEnv(&Global); err != nil {
		log.Fatalf("config: can't read envs: %v", err)
	}
}
