package config

import (
	"github.com/go-yaml/yaml"
	"log"
	"os"
)

type Config struct {
	Server struct {
		Name     string   `yaml:"name"`
		Host     string   `yaml:"host"`
		Port     string   `yaml:"port"`
		Comrades []string `yaml:"comrades"`
	}
	Logger struct {
		Type     string `yaml:"type"`
		Path     string `yaml:"path"`
		Rotating bool   `yaml:"rotating"`
	}
	Watcher struct {
		Directories map[string]string `yaml:"directories"`
	}
}

type ErrorConfig struct {
}

func (e *ErrorConfig) Error() string {
	return "Config not provided. Be sure that config/config.yaml exists"
}

func newErrorConfig() *ErrorConfig {
	return &ErrorConfig{}
}

func attachDefaults(cfg *Config) {
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if cfg.Server.Host == "" {
		cfg.Server.Host = "localhost"
	}
	if cfg.Server.Name == "" {
		cfg.Server.Name = "localhost"
	}
	if cfg.Logger.Type == "" {
		cfg.Logger.Type = "console"
	}
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, newErrorConfig()
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	yamlDecoder := yaml.NewDecoder(file)
	if err := yamlDecoder.Decode(&config); err != nil {
		log.Fatal(err)
		return nil, err
	}
	attachDefaults(config)
	return config, nil
}
