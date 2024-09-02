/*
Copyright © 2024 Stefan Victor (faneatiku@yahoo.com)
*/
package main

import (
	log "github.com/sirupsen/logrus"
	"go-tha-utxos/app/action"
	"go-tha-utxos/cmd"
	"os"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer
	log.SetOutput(os.Stdout)
}

func main() {
	// Explicitly check if there are no command-line arguments
	if len(os.Args) == 1 {
		// No arguments provided, assume double-click or direct execution
		err := action.RunApp()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Execute the root command with arguments (normal CLI behavior)
		cmd.Execute()
	}
}
