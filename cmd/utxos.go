/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
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
}
