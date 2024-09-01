package action

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-tha-utxos/app/services"
	"go-tha-utxos/config"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunApp() error {
	cfg, err := config.LoadAndApplyConfig()
	if err != nil {
		return fmt.Errorf("could not load config: %v", err)
	}

	err = createNeededAddresses(cfg)
	if err != nil {
		return fmt.Errorf("error while creating needed addresses: %v", err)
	}

	quit := make(chan struct{})
	addSigtermHandler(quit)

	go runRecursive(cfg)

	<-quit

	return nil
}

func runRecursive(cfg *config.Config) {
	log.Info(fmt.Sprintf("will try to generate UTXOs in %d minutes", cfg.AutoRunner.UtxosInterval))
	time.Sleep(time.Duration(cfg.AutoRunner.UtxosInterval) * time.Minute)

	err := CreateUtxos(cfg, cfg.AutoRunner.AddressesFile, cfg.AutoRunner.UtxosFee)
	if err != nil {
		log.WithError(err).Error("create utxos failed")
	}

	runRecursive(cfg)
}

func createNeededAddresses(cfg *config.Config) error {
	l := log.WithField("addresses_file", cfg.AutoRunner.AddressesFile).
		WithField("addresses_count", cfg.AutoRunner.AddressesCount)

	if services.FileExists(cfg.AutoRunner.AddressesFile) {
		l.Info("addresses file already exists")
		l.Info("using already existing addresses")

		return nil
	}

	l.Info("creating addresses file to use for UTXOs")

	return CollectAddresses(cfg, cfg.AutoRunner.AddressesCount, cfg.AutoRunner.AddressesFile, false)
}

func addSigtermHandler(quit chan struct{}) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Info("sigterm received. shutting down")
		quit <- struct{}{}
		os.Exit(0)
	}()
}
