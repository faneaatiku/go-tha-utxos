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

// collectAddressesCmd represents the collectAddresses command
var collectAddressesCmd = &cobra.Command{
	Use:   "collect",
	Short: "Collects all addresses of the node",
	Long: `Collects all addresses of a loaded wallet using listaddressgroupings cli command and shows them on the screen
or saves them to a json file. The number of addresses needed can be passed and the command will generate generate new addresses
if it doesn't find enough already generated.

Usage:
1. Collect all addresses found: 
./go-tha-utxos addresses collect

2. Collect 100 addresses - if not enough addresses found it will generate new ones in order to return the 
number of addresses requested:
./go-tha-utxos addresses collect --count 100

2. Collect 100 addresses and save them to a json file called file_name.json - if not enough addresses found it will generate new ones in order to return the 
number of addresses requested:
./go-tha-utxos addresses collect  --count 100 --file file_name.json
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		count, err := cmd.Flags().GetInt(countFlag)
		if err != nil {
			return err
		}

		if count < 0 || count > maxCount {
			return fmt.Errorf("the count must be from %d to %d", 0, maxCount)
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

		return action.CollectAddresses(cfg, count, file, ignoreExistingFile)
	},
}

func init() {
	addressesCmd.AddCommand(collectAddressesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// collectAddressesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// collectAddressesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
