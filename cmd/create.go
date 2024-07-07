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
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	feeFlag   = "fee"
	keepAlive = "keep-alive"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Command to create minable utxos",
	Long: `This command searches for UTXOs that can be used to create minable UTXOs of value 0.1

--file flag (mandatory) the json file containing the addresses to send UTXOs to. The file should be generated with "addresses generate" or "addresses collect" command
--fee flag (optional) to specify the fee rate in sats/byte (default 0 which will let the command calculate the fee)

Example: 
1. Create at least 100 UTXOs - let the command decide the TX fee to use
./go-tha-utxos utxos create --count 100

2. Create at least 100 UTXOs with 15 sats/byte fee
./go-tha-utxos utxos create --count 100 --fee 15
`,
	RunE: func(cmd *cobra.Command, args []string) error {

		fee, err := cmd.Flags().GetFloat64(feeFlag)
		if err != nil {
			return err
		}

		file, err := cmd.Flags().GetString(fileFlag)
		if err != nil {
			return err
		}

		if file == "" {
			log.Fatal("please provide a file that contains addresses for UTXOs using --file flag")
		}

		cfg, err := config.LoadAndApplyConfig()
		if err != nil {
			return fmt.Errorf("could not load config: %v", err)
		}

		keepAlive, err := cmd.Flags().GetBool(keepAlive)
		if err != nil {
			log.Errorf("keep alive flag could not be used: %v", err)
			keepAlive = false
		}

		result := action.CreateUtxos(cfg, file, fee)
		if !keepAlive {
			return result
		}

		ticker := time.NewTicker(2 * time.Minute)
		quit := make(chan struct{})
		addSigtermHandler(quit)
		for {
			select {
			case <-ticker.C:
				log.Info("running create utxos")
				err := action.CreateUtxos(cfg, file, fee)
				if err != nil {
					log.Errorf("create utxos flow returned error: %v", err)
				} else {
					log.Info("create utxos successfully finished")
				}
			case <-quit:
				ticker.Stop()
				return nil
			}
		}
	},
}

func init() {
	utxosCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().Float64(feeFlag, 0.01, "the fee to use for the transaction")
	createCmd.Flags().String(fileFlag, "", "the file with the addresses to send UTXOs to")
	createCmd.Flags().Bool(keepAlive, false, "if passed this flag will keep the command alive and will run the action every 2 minutes")
}

func addSigtermHandler(quitChan chan struct{}) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(quitChan)
	}()
}
