package config

import log "github.com/sirupsen/logrus"

type Config struct {
	Commands Commands `yaml:"commands"`
	Logging  Logging  `yaml:"logging"`
}

type Logging struct {
	Level       string    `yaml:"level"`
	ParsedLevel log.Level `yaml:"-"`
}

type Commands struct {
	DaemonCli string `yaml:"daemon_cli"`
}
