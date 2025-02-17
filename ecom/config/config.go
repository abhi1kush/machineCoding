package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config holds the configuration settings from the YAML file.
type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Driver string `yaml:"driver"`
		DSN    string `yaml:"dsn"`
	} `yaml:"database"`
	Metrics struct {
		Driver string `yaml:"driver"`
		DSN    string `yaml:"dsn"`
	} `yaml:"metrics"`
	Queue struct {
		WorkerPool    int `yaml:"workerPool"`
		QueueCapacity int `yaml:"queueCapacity"`
	} `yaml:"queue"`
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
}

// AppConfig is a global instance of Config.
var AppConfig Config

// LoadConfig reads configuration from the specified YAML file.
func LoadConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
}
