package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
		PG   `yaml:"postgres"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		Host          string `env-required:"true" yaml:"host" env:"PG_HOST"`
		Port          string `env-required:"true" yaml:"port" env:"PG_PORT"`
		Name          string `env-required:"true" yaml:"name" env:"PG_NAME"`
		User          string `env-required:"true" yaml:"user" env:"PG_USER"`
		Password      string `env-required:"true" yaml:"password" env:"PG_PASSWORD"`
		SSLMode       string `env-required:"true" yaml:"ssl_mode" env:"PG_SSL_MODE"`
		LogLevel      int    `yaml:"log_level" env:"PG_LOG_LEVEL"`
		AutoMigrate   bool   `yaml:"auto_migrate" env:"PG_AUTO_MIGRATE"`
		GenerateSeeds bool   `yaml:"generate_seeds" env:"PG_GENERATE_SEEDS"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
