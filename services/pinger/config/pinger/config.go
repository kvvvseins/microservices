package pinger

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
		Connection    *gorm.DB `ignored:"true"`
		TYPE          string   `envconfig:"DB_TYPE" required:"true"`
		ConnectionDsn string   `envconfig:"CONNECTION_DSN" required:"true"`
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

	dbConnection(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
	defer cancel()

	sqlDb, err := cfg.DB.Connection.DB()
	if err != nil {
		log.Fatalf("failed to get connect to db: %v", err)
	}

	err = sqlDb.PingContext(ctx)
	if err != nil {
		log.Fatalf("failed to ping to db: %v", err)
	}

	return nil
}

func dbConnection(cfg *Config) {
	var err error

	switch cfg.DB.TYPE {
	case "postgres":
		cfg.DB.Connection, err = gorm.Open(postgres.Open(cfg.DB.ConnectionDsn), &gorm.Config{})
	default:
		log.Fatalf("connect type not supported: %s", cfg.DB.TYPE)
	}

	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
}
