package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string `yaml:"env" env:"ENV" env-default:"local"`
	LimitCache     int    `yaml:"limitC" env-default:"50"`
	MaxItemsCache  int    `yaml:"maxItemsC" env-default:"100"`
	HTTPServer     `yaml:"http_server"`
	PostgresConfig `yaml:"db"`
	KafkaConfig    `yaml:"kafka"`
}

type HTTPServer struct {
	Port        string        `yaml:"port" env-default:"8082"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type PostgresConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type KafkaConfig struct {
	Brokers    []string      `yaml:"brokers"`
	Topic      string        `yaml:"topic"`
	GroupID    string        `yaml:"group_id"`
	DLQTopic   string        `yaml:"dlq_topic"`
	RetryDelay time.Duration `yaml:"retry_delay"`
}

func Load() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("CONFIG_PATH is not set")
	}

	//check if file exists
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file does not exist: %s", configPath)
		}
		return nil, fmt.Errorf("cannot stat config file %q: %w", configPath, err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	return &cfg, nil

}
