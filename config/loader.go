package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	defaultLoggingLevel = log.InfoLevel

	defaultAddressesFile       = "auto.addresses.json"
	defaultAddressesCount      = 500
	defaultUtxosInterval       = 2
	defaultUtxosFee            = 0.001
	defaultConsolidateMinUtxos = 100
	defaultConsolidateInterval = 1440 //1 day
)

var config *Config

func LoadConfig() (*Config, error) {
	if config != nil {
		return config, nil
	}

	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config.yaml: %w", err)
	}

	cfg := &Config{}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config.yaml: %w", err)
	}

	if cfg.Logging.Level != "" {
		level, err := log.ParseLevel(cfg.Logging.Level)
		if err != nil {
			return nil, fmt.Errorf("failed to parse logging level from config.yaml: %w", err)
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

	if cfg.AutoRunner.AddressesFile == "" {
		cfg.AutoRunner.AddressesFile = defaultAddressesFile
	}

	if cfg.AutoRunner.AddressesCount == 0 {
		cfg.AutoRunner.AddressesCount = defaultAddressesCount
	}

	if cfg.AutoRunner.UtxosInterval <= 0 {
		cfg.AutoRunner.UtxosInterval = defaultUtxosInterval
	}

	if cfg.AutoRunner.UtxosFee <= 0 {
		cfg.AutoRunner.UtxosFee = defaultUtxosFee
	}

	if cfg.AutoRunner.UtxosFee > 1 {
		log.Fatal("utxos_fee is too high! Recommended value 0.01")
	}

	if cfg.AutoRunner.ConsolidateMinUtxos <= 0 {
		cfg.AutoRunner.UtxosFee = defaultConsolidateMinUtxos
	}

	if cfg.AutoRunner.ConsolidateInterval <= 0 {
		cfg.AutoRunner.ConsolidateInterval = defaultConsolidateInterval
	}
}
