package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
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
	Host                string        `yaml:"host"`
	Port                string        `yaml:"port"`
	Base                string        `yaml:"base"`
	User                string        `yaml:"user"`
	Password            string        `yaml:"password"`
	MinConns            int32         `yaml:"min_conns"`
	MaxConns            int32         `yaml:"max_conns"`
	MaxLife             time.Duration `yaml:"max_life"`
	MaxIdle             time.Duration `yaml:"max_idle"`
	HealthCheckInterval time.Duration `yaml:"health_check_interval"`
}

type RedisConfig struct {
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	Base        string        `yaml:"base"`
	Password    string        `yaml:"password"`
	MaxConn     int           `yaml:"max_conn"`
	MaxIdle     int           `yaml:"max_idle"`
	MaxActive   int           `yaml:"max_active"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func NewConfig() (*Config, error) {
	configPath := "config/config.yaml"

	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config appConfig
	config.App.Postgres = PostgresConfig{
		MinConns:            1,
		MaxConns:            5,
		MaxLife:             time.Hour,
		MaxIdle:             30 * time.Minute,
		HealthCheckInterval: time.Minute,
	}
	config.App.Redis = RedisConfig{
		MaxConn:     5,
		MaxIdle:     10,
		MaxActive:   0,
		IdleTimeout: 4 * time.Minute,
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config.App, nil
}
