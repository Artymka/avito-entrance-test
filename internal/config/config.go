package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Postgres struct {
		Host     string `yaml:"host" validate:"required"`
		Port     string `yaml:"port" validate:"required"`
		Database string `yaml:"database" validate:"required"`
		SSLMode  string `yaml:"ssl_mode" validate:"required"`
		Username string
		Password string
	} `yaml:"postgres"`
}

func New(configPath string) (*Config, error) {
	const op = "config.new"
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	cfg := Config{}
	if err = yaml.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	v := validator.New()
	if err = v.Struct(cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// env vars
	cfg.Postgres.Username = os.Getenv("POSTGRES_USER")
	if cfg.Postgres.Username == "" {
		cfg.Postgres.Username = "root"
	}

	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	if cfg.Postgres.Password == "" {
		cfg.Postgres.Password = "root"
	}

	return &cfg, nil
}
