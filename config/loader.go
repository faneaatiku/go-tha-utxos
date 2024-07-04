package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	defaultLoggingLevel = log.InfoLevel
)

var config *Config

func LoadConfig() (*Config, error) {
	if config != nil {
		return config, nil
	}

	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config.yml: %w", err)
	}

	cfg := &Config{}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config.yml: %w", err)
	}

	if cfg.Logging.Level != "" {
		level, err := log.ParseLevel(cfg.Logging.Level)
		if err != nil {
			return nil, fmt.Errorf("failed to parse logging level from config.yml: %w", err)
		}

		cfg.Logging.ParsedLevel = level
	} else {
		cfg.Logging.ParsedLevel = defaultLoggingLevel
	}

	config = cfg

	return config, nil
}

func LoadAndApplyConfig() (*Config, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	applyGlobalConfig(cfg)

	return cfg, nil
}

func applyGlobalConfig(cfg *Config) {
	if cfg == nil {
		return
	}

	// Only log the warning severity or above.
	log.SetLevel(cfg.Logging.ParsedLevel)
}
