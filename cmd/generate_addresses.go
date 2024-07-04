/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-tha-utxos/app/action"
	"go-tha-utxos/config"
)

const (
	minCount = 1
	maxCount = 10000

	countFlag = "count"
)

// generateAddressesCmd represents the generateAddresses command
var generateAddressesCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a defined number of addresses and saves them to a csv file",
	Long: `Generate addresses in the loaded wallet using getnewaddress cli command and shows them on the screen
or saves them to a json file. The number of addresses needed can be passed and the command will generate generate new addresses
if it doesn't find enough already generated.

Usage:
1. Create 10 new addresses and show them on screen: 
./go-tha-utxos addresses generate --count 10

2. Create 100 addresses and save them to addresses.json:
./go-tha-utxos addresses generate --count 100 --file addresses 
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		count, err := cmd.Flags().GetInt(countFlag)
		if err != nil {
			return err
		}

		if count < minCount || count > maxCount {
			return fmt.Errorf("the --count flag must have values between %d and %d", minCount, maxCount)
		}

		cfg, err := config.LoadAndApplyConfig()
		if err != nil {
			return fmt.Errorf("could not load config: %v", err)
		}

		file, _ := cmd.Flags().GetString(fileFlag)
		if file == "" {
			log.Debug("no file specified for generate addresses")
		}
		ignoreExistingFile := parseIgnoreExistingFile(cmd)

		return action.GenerateAddresses(cfg, count, file, ignoreExistingFile)
	},
}

func init() {
	addressesCmd.AddCommand(generateAddressesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateAddressesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//generateAddressesCmd.Flags().String(fileFlag, "", "the file in which to write the addresses")
}
