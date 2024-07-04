/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	fileFlag           = "file"
	ignoreExistingFile = "ignore-existing-file"
)

// addressesCmd represents the addresses command
var addressesCmd = &cobra.Command{
	Use:   "addresses",
	Short: "Group of commands related to addresses",
	Long: `
collect -- collect existing addresses from daemon and add new ones if needed
generate -- generate a fix number of addresses from the daemon
`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Usage()
	},
}

func init() {
	rootCmd.AddCommand(addressesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addressesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addressesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	addressesCmd.PersistentFlags().String(fileFlag, "", "the file in which to write the addresses")
	addressesCmd.PersistentFlags().Bool(ignoreExistingFile, false, "override the provide file if it exists")
	addressesCmd.PersistentFlags().Int(countFlag, 0, "the number of addresses to collect. Use 0 to collect all of them and not create new ones")
}

func parseIgnoreExistingFile(cmd *cobra.Command) bool {
	ignoreFlag, err := cmd.Flags().GetBool(ignoreExistingFile)
	if err != nil {
		log.Warnf("failed to parse ignore-existing-file flag: %v. Using default false", err)
		ignoreFlag = false
	}

	return ignoreFlag
}
