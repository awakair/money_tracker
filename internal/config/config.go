package config

import (
	_ "embed"
	"fmt"
	"os"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"gopkg.in/yaml.v3"
)

//go:embed default_config.yaml
var defaultConfigData []byte

type Config struct {
	Database databaseConfig `yaml:"database"`
	Server   serverConfig   `yaml:"server"`
	JWT      jwtConfig      `yaml:"jwt"`
	Worker   workerConfig   `yaml:"worker"`
	Log      logConfig      `yaml:"log"`
}

type databaseConfig struct {
	URL string `yaml:"url"`
}

func (d databaseConfig) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.URL, validation.Required),
	)
}

type serverConfig struct {
	Port string `yaml:"port"`
}

func (s serverConfig) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Port, validation.Required),
	)
}

type jwtConfig struct {
	Secret     string        `yaml:"secret"`
	Expiration time.Duration `yaml:"expiration"`
}

func (j jwtConfig) Validate() error {
	return validation.ValidateStruct(&j,
		validation.Field(&j.Secret, validation.Required),
		validation.Field(&j.Expiration, validation.Required, validation.Min(time.Second)),
	)
}

type workerConfig struct {
	PollInterval time.Duration `yaml:"poll_interval"`
}

func (w workerConfig) Validate() error {
	return validation.ValidateStruct(&w,
		validation.Field(&w.PollInterval, validation.Required, validation.Min(time.Second)),
	)
}

type logConfig struct {
	Level string `yaml:"level"`
}

func (l logConfig) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Level, validation.Required, validation.In("debug", "info", "warn", "error")),
	)
}

func LoadFromFile(filepath string) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(defaultConfigData, &cfg); err != nil {
		return nil, fmt.Errorf("parse default config: %w", err)
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Database),
		validation.Field(&c.Server),
		validation.Field(&c.JWT),
		validation.Field(&c.Worker),
		validation.Field(&c.Log),
	)
}
