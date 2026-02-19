package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpiryHours int    `yaml:"expiry_hours"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Database.Path == "" {
		cfg.Database.Path = "./data/web.db"
	}
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "go-etl-web-secret-key"
	}
	if cfg.JWT.ExpiryHours == 0 {
		cfg.JWT.ExpiryHours = 24
	}

	return &cfg, nil
}
