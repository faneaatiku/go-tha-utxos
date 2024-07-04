/*
Copyright © 2024 Stefan Victor (faneatiku@yahoo.com)
*/
package main

import (
	log "github.com/sirupsen/logrus"
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
	cmd.Execute()
}
