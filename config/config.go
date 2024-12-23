package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

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
}

// LoadConfig функция для загрузки конфигурации.
func LoadConfig(cfg interface{}, prefix string) error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return errors.Wrap(err, "не удалось создать конфиг")
	}

	if err := envconfig.Process(prefix, cfg); err != nil {
		return errors.Wrap(err, "ошибка обработки env-переменных")
	}

	return nil
}
