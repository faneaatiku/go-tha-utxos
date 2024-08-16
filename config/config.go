package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Commands      Commands      `yaml:"commands"`
	Logging       Logging       `yaml:"logging"`
	RpcConnection RpcConnection `yaml:"rpc"`
}

type Logging struct {
	Level       string    `yaml:"level"`
	ParsedLevel log.Level `yaml:"-"`
}

type Commands struct {
	DaemonCli string `yaml:"daemon_cli"`
	DataDir   string `yaml:"data_dir"`
}

type RpcConnection struct {
	Host       string `yaml:"host"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	WalletName string `yaml:"wallet_name"`
}

func (c RpcConnection) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("rpc host is required")
	}

	if c.User == "" {
		return fmt.Errorf("rpc user is required")
	}

	if c.Password == "" {
		return fmt.Errorf("rpc password is required")
	}

	return nil
}
