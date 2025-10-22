package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type appConfig struct {
	App Config `yaml:"app"`
}

type Config struct {
	Host     string         `yaml:"host"`
	Port     int            `yaml:"port"`
	Prefix   string         `yaml:"prefix"`
	Logger   LoggerConfig   `yaml:"logger"`
	Postgres PostgresConfig `yaml:"postgres"`
	Redis    RedisConfig    `yaml:"redis"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Base     string `yaml:"base"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Base     string `yaml:"base"`
	Password string `yaml:"password"`
}

func NewConfig() (*Config, error) {
	configPath := "config/config.yaml"

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config appConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config.App, nil
}
