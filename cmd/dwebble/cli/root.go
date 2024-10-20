package cli

import (
	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     dwebble.AppName,
	Version: dwebble.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Failed to run root command", "err", err)
	}
}

func init() {
	rootCmd.SetVersionTemplate(dwebble.VersionTemplate(dwebble.AppName))

	cobra.OnInitialize(config.InitConfig)

	rootCmd.PersistentFlags().StringVarP(&config.ConfigFile, "config", "c", "", "Config File")
}
