package config

import (
	"fmt"
	"github.com/caarlos0/env/v10"
)

type Config struct {
	AppPort string `env:"APP_PORT" envDefault:"8080"`

	DBHost     string `env:"DB_HOST"     envDefault:"localhost"`
	DBPort     string `env:"DB_PORT"     envDefault:"5432"`
	DBUser     string `env:"DB_USER"     envDefault:"postgres"`
	DBPassword string `env:"DB_PASS"     envDefault:"postgres"`
	DBName     string `env:"DB_NAME"     envDefault:"app_db"`
	SSLMode    string `env:"DB_SSLMODE"  envDefault:"disable"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config load error: %w", err)
	}
	return cfg, nil
}

func (c *Config) BuildDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.SSLMode,
	)
}
