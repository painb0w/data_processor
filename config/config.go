package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type Config struct {
	FootballApi FootballApiConfig `yaml:"football-api"`
	Redis       RedisConfig       `yaml:"redis"`
	Mongo       MongoDBConfig     `yaml:"mongodb"`
	Kafka       KafkaConfig       `yaml:"kafka"`
}

type FootballApiConfig struct {
	ApiKey  string `yaml:"api_key"`
	BaseUrl string `yaml:"base_url"`
	Timeout int    `yaml:"timeout"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
}

type MongoDBConfig struct {
	Uri string `yaml:"uri"`
}

type KafkaConfig struct {
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	SubscriptionTopic string `yaml:"subscription_topic"`
	StatsUpdateTopic  string `yaml:"stats_update_topic"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &config, nil
}
