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

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export addresses from a file and writes them to a file",
	Long:  `Export addresses from the provided file to a new file`,
	RunE: func(cmd *cobra.Command, args []string) error {

		cfg, err := config.LoadAndApplyConfig()
		if err != nil {
			return fmt.Errorf("could not load config: %v", err)
		}

		file, _ := cmd.Flags().GetString(fileFlag)
		if file == "" {
			log.Debug("no file specified for generate addresses")
		}

		return action.ExportAddresses(cfg, file)
	},
}

func init() {
	addressesCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
