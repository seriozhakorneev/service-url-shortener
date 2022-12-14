package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App       `yaml:"app"`
		HTTP      `yaml:"http"`
		Log       `yaml:"logger"`
		PG        `yaml:"postgres"`
		Redis     `yaml:"redis"`
		GRPC      `yaml:"grpc"`
		Digitiser `yaml:"digitiser"`
		URL       `yaml:"url"`
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
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true"                 env:"PG_URL"`
	}

	// Redis -.
	Redis struct {
		Address string `env-required:"true" env:"REDIS_ADDRESS"`
		Pass    string `env-required:"true" env:"REDIS_PASS"`
		DB      int    `yaml:"db"           env:"REDIS_DB"`
	}

	// GRPC -.
	GRPC struct {
		Network string `env-required:"true" yaml:"network" env:"GRPC_NETWORK"`
		Port    string `env-required:"true" yaml:"port" env:"GRPC_PORT"`
	}

	// Digitiser -.
	Digitiser struct {
		MaxCount  int    `yaml:"max_count"`
		Base      string `yaml:"base"`
		MaxLength int    `yaml:"max_length"`
	}

	// URL -.
	URL struct {
		Blank string `yaml:"blank"`
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
