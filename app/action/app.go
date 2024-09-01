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
	l := log.WithField("action", "startup")

	l.Info(fmt.Sprintf("will try to create new UTXOs every %d minutes", cfg.AutoRunner.UtxosInterval))
	go runCreateUtxos(cfg)

	l.Info(fmt.Sprintf("will try to consolidate UTXOs every %d minutes", cfg.AutoRunner.ConsolidateInterval))
	go runConsolidateUtxos(cfg)

	<-quit

	return nil
}

func runCreateUtxos(cfg *config.Config) {
	l := log.WithField("action", "utxos/create")
	time.Sleep(time.Duration(cfg.AutoRunner.UtxosInterval) * time.Minute)

	l.Info("try to create new UTXOs")
	err := CreateUtxos(cfg, cfg.AutoRunner.AddressesFile, cfg.AutoRunner.UtxosFee)
	if err != nil {
		l.WithError(err).Error("create UTXOs failed")
	}

	runCreateUtxos(cfg)
}

func runConsolidateUtxos(cfg *config.Config) {
	l := log.WithField("action", "utxos/consolidate")
	time.Sleep(time.Duration(cfg.AutoRunner.ConsolidateInterval) * time.Minute)

	l.Info("try to consolidate UTXOs")
	err := ConsolidateUtxos(cfg, cfg.AutoRunner.UtxosFee, cfg.AutoRunner.ConsolidateMinUtxos)
	if err != nil {
		l.WithError(err).Error("consolidate UTXOs failed")
	}

	runConsolidateUtxos(cfg)
}

func createNeededAddresses(cfg *config.Config) error {
	l := log.WithField("addresses_file", cfg.AutoRunner.AddressesFile).
		WithField("addresses_count", cfg.AutoRunner.AddressesCount).
		WithField("action", "startup")

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
