package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type appConfig struct {
	App Config `yaml:"app"`
}

type Config struct {
	Host        string            `yaml:"host"`
	Port        int               `yaml:"port"`
	Prefix      string            `yaml:"prefix"`
	Cors        CORSConfig        `yaml:"cors"`
	Logger      LoggerConfig      `yaml:"logger"`
	Postgres    PostgresConfig    `yaml:"postgres"`
	Redis       RedisConfig       `yaml:"redis"`
	MinIO       MinIOConfig       `yaml:"minio"`
	AuthService AuthServiceConfig `yaml:"auth_service"`
}

type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAgeSeconds    int      `yaml:"max_age_seconds"`
}

type LoggerConfig struct {
	Level     string `yaml:"level"`
	Prefix    string `yaml:"prefix"`
	Color     bool   `yaml:"color"`
	Timestamp bool   `yaml:"timestamp"`
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

	Migrated   bool   `yaml:"migrated"`
	Migrations string `yaml:"migrations"`
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

type MinIOConfig struct {
	Host            string `yaml:"host"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	Address         string `yaml:"address"`
	UseSSL          bool   `yaml:"use_ssl"`
	BucketName      string `yaml:"bucket_name"`
	Region          string `yaml:"region"`
}

type AuthServiceConfig struct {
	Port string `yaml:"port"`
}

func NewConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config appConfig
	config.App.Cors = CORSConfig{
		AllowedOrigins:   []string{},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With", "Accept", "Origin"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAgeSeconds:    86400,
	}
	config.App.Logger = LoggerConfig{
		Level:     "INFO",
		Prefix:    "",
		Color:     true,
		Timestamp: true,
	}
	config.App.Postgres = PostgresConfig{
		MinConns:            1,
		MaxConns:            5,
		MaxLife:             time.Hour,
		MaxIdle:             30 * time.Minute,
		HealthCheckInterval: time.Minute,

		Migrated:   false,
		Migrations: "migrations/",
	}
	config.App.Redis = RedisConfig{
		MaxConn:     5,
		MaxIdle:     10,
		MaxActive:   0,
		IdleTimeout: 240 * time.Second,
	}
	config.App.MinIO = MinIOConfig{
		UseSSL: false,
		Region: "us-east-1",
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config.App, nil
}
