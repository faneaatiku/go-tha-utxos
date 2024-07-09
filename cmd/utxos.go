/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	feeFlag      = "fee"
	intervalFlag = "interval"
)

// utxosCmd represents the utxos command
var utxosCmd = &cobra.Command{
	Use:   "utxos",
	Short: "UTXOs commands",
	Long: `Commands related to UTXOs :
1. Create UTXOs command 
./go-tha-utxos utxos create --count 10
`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Usage()
	},
}

func init() {
	rootCmd.AddCommand(utxosCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// utxosCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// utxosCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	utxosCmd.PersistentFlags().Float64(feeFlag, 0.001, "the fee to use for the transaction. Default 0.001")
	utxosCmd.PersistentFlags().Int(intervalFlag, 0, "interval in minutes. If passed the command will run every X minutes. if 0 the command will run once. Default 0")
}

func mustGetIntervalFlag(cmd *cobra.Command) int {
	intervalFlagValue, err := cmd.Flags().GetInt(intervalFlag)
	if err != nil {
		log.Errorf("%s flag could not be used: %v", intervalFlag, err)
		intervalFlagValue = 0
	}

	return intervalFlagValue
}

func runDelayed(cmd *cobra.Command, delayedFunc func() error) error {
	interval := mustGetIntervalFlag(cmd)

	result := delayedFunc()
	if interval <= 0 {
		return result
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Minute)
	quit := make(chan struct{})
	addSigtermHandler(quit)
	for {
		select {
		case <-ticker.C:
			err := delayedFunc()
			if err != nil {
				log.Errorf("delayed function error: %v", err)
			}
		case <-quit:
			ticker.Stop()
			return nil
		}
	}
}

func addSigtermHandler(quitChan chan struct{}) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(quitChan)
	}()
}
