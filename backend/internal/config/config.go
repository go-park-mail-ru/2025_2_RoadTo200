package config

import (
	"fmt"
	"log"
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
		log.Printf("⚠️  No .env file found: %v", err)
	} else {
		log.Println("✅ .env file loaded")
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	log.Printf("✅ Config file read: %s", configPath)

	var config appConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	log.Printf("🔍 Config after YAML parse: Host=%s, Port=%d", config.App.Host, config.App.Port)
	log.Printf("🔍 PostgreSQL Config: Host=%s, Port=%s, User=%s",
		config.App.Postgres.Host, config.App.Postgres.Port, config.App.Postgres.User)
	log.Printf("🔍 Redis Config: Host=%s, Port=%s",
		config.App.Redis.Host, config.App.Redis.Port)

	return &config.App, nil
}
