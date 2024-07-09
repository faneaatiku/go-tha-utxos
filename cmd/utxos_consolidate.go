/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-tha-utxos/app/action"
	"go-tha-utxos/config"

	"github.com/spf13/cobra"
)

const (
	minUtxosFlag = "min-utxos"

	defaultMinUtxos = 50
)

// consolidateCmd represents the consolidate command
var consolidateCmd = &cobra.Command{
	Use:   "consolidate",
	Short: "Consolidates amount bigger than 0.1 and lower than 0.3",
	Long:  `Consolidates amount bigger than 0.1 and lower than 0.3`,
	RunE: func(cmd *cobra.Command, args []string) error {

		fee, err := cmd.Flags().GetFloat64(feeFlag)
		if err != nil {
			return err
		}

		minUtxos, err := cmd.Flags().GetInt(minUtxosFlag)
		if err != nil {
			log.Errorf("could not parse %s flag due to error: %v. Using default :%d", minUtxosFlag, err, defaultMinUtxos)
			minUtxos = defaultMinUtxos
		}

		if minUtxos <= 5 {
			return fmt.Errorf("%s flag must be greater than 4", minUtxosFlag)
		}

		cfg, err := config.LoadAndApplyConfig()
		if err != nil {
			return fmt.Errorf("could not load config: %v", err)
		}

		delayedFunc := func() error {
			log.Info("running consolidate utxos")
			err := action.ConsolidateUtxos(cfg, fee, minUtxos)
			if err != nil {
				log.Errorf("consolidate utxos flow returned error: %v", err)
			} else {
				log.Info("consolidate utxos successfully finished")
			}

			return err
		}

		return runDelayed(cmd, delayedFunc)
	},
}

func init() {
	utxosCmd.AddCommand(consolidateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consolidateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consolidateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	consolidateCmd.PersistentFlags().Int(minUtxosFlag, defaultMinUtxos, fmt.Sprintf("The minimum number of extra UTXOs that should result after consolidation. Default %d", defaultMinUtxos))
}
