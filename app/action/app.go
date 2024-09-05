package action

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-tha-utxos/app/dto/daemon"
	"go-tha-utxos/app/math"
	"go-tha-utxos/app/services"
	"go-tha-utxos/config"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"time"
)

var d *services.RpcDaemon

func RunApp() error {
	cfg, err := config.LoadAndApplyConfig()
	if err != nil {
		return fmt.Errorf("could not load config: %v", err)
	}

	d, err = services.NewRpcDaemon(&cfg.RpcConnection)
	if err != nil {
		return fmt.Errorf("could not create RPC daemon: %v", err)
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

	if isMaxHashRateReached(cfg) {
		l.Info("max hash rate reached")
		if cfg.AutoRunner.ExtraHashrateAddress != "" {
			l.Info("checking for eligible funds that are not minable UTXOs to send to extra_hashrate_address")
			err := exportFunds(l, cfg)
			if err != nil {
				l.WithError(err).Error("could not export funds that are not minable UTXOs to send to extra_hashrate_address")
			}
		} else {
			l.Error("can not send extra funds exceeding the max_hashrate to another address: no extra_hashrate_address provided")
		}
	} else {
		l.Info("try to create UTXOs")
		err := CreateUtxos(cfg, cfg.AutoRunner.AddressesFile, cfg.AutoRunner.UtxosFee)
		if err != nil {
			l.WithError(err).Error("create UTXOs failed")
		}
	}

	runCreateUtxos(cfg)
}

func exportFunds(l *log.Entry, cfg *config.Config) error {
	unspent, err := d.ListUnspent(500)
	if err != nil {
		return fmt.Errorf("finding unspent addresses failed: %s", err)
	}

	if len(unspent) == 0 {
		l.Info("no unspent eligible funds found")
		return nil
	}

	slices.SortStableFunc(unspent, func(a, b daemon.Unspent) int {
		return int(b.Amount) - int(a.Amount)
	})

	feeDec, err := math.NewDecFromStr("0.00001")
	if err != nil {
		return fmt.Errorf("error on fee conversion 0.01: %s", err)
	}
	//increase the fee to make sure it's enough
	feeDec = feeDec.MulInt64(int64(len(unspent)))

	amountFound := math.ZeroDec()
	var selectedUnspent []daemon.RawTransactionInput
	for i := 0; i < len(unspent); i++ {
		unspendAmount, err := math.NewDecFromStr(services.FloatToString(unspent[i].Amount))
		amountFound = amountFound.Add(unspendAmount)
		if err != nil {
			return fmt.Errorf("error on unspend amount conversion [%.2f]: %s", unspendAmount, err)
		}
		selectedUnspent = append(selectedUnspent, daemon.RawTransactionInput{
			Txid: unspent[i].Txid,
			Vout: unspent[i].Vout,
		})
	}

	var outputs []daemon.RawTransactionOutput
	raw := make(daemon.RawTransactionOutput, 1)
	raw[cfg.AutoRunner.ExtraHashrateAddress] = amountFound.Sub(feeDec).MustFloat64()
	outputs = append(outputs, raw)

	sendTransaction(d, selectedUnspent, outputs)

	return nil
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

func isMaxHashRateReached(cfg *config.Config) bool {
	mi, err := d.GetMiningInfo()
	if err != nil {
		log.WithError(err).Error("error getting mining info")

		return false
	}

	return mi.Localhashps >= cfg.AutoRunner.MaxHashrate
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
