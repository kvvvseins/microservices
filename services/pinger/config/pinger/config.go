package pinger

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

import _ "github.com/lib/pq"

// Config конфигурации приложения.
type Config struct {
	App struct {
		Name     string     `envconfig:"NAME" required:"true"`
		LogLevel slog.Level `default:"warn" envconfig:"LOG_LEVEL"`
	}
	Env  string `default:"dev" envconfig:"ENV"`
	HTTP struct {
		Port int `default:"8000" envconfig:"PORT"`
	}
	DB struct {
		Connection    *sql.DB
		TYPE          string `default:"" envconfig:"DB_TYPE"`
		ConnectionDsn string `default:"" envconfig:"CONNECTION_DSN"`
	}
}

// LoadConfig функция для загрузки конфигурации.
func LoadConfig(cfg *Config, prefix string) error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return errors.Wrap(err, "не удалось создать конфиг")
	}

	if err := envconfig.Process(prefix, cfg); err != nil {
		return errors.Wrap(err, "ошибка обработки env-переменных")
	}

	var err error
	cfg.DB.Connection, err = sql.Open(cfg.DB.TYPE, cfg.DB.ConnectionDsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
	defer cancel()

	err = cfg.DB.Connection.PingContext(ctx)
	fmt.Println(err)
	if err != nil {
		log.Fatalf("failed to ping to db: %v", err)
	}
	fmt.Println(cfg.DB.ConnectionDsn)

	return nil
}
