package cmd

import (
	"log"

	"github.com/nanoteck137/dwebble"
	"github.com/spf13/cobra"
)

var AppName = dwebble.AppName + "-migrate"

var rootCmd = &cobra.Command{
	Use:     AppName,
	Version: dwebble.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.SetVersionTemplate(dwebble.VersionTemplate(AppName))
}
