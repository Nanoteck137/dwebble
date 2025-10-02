package cmd

import (
	"log/slog"
	"os"

	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     dwebble.AppName,
	Version: dwebble.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Failed to run root command", "err", err)
		os.Exit(-1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(dwebble.VersionTemplate(dwebble.AppName))

	cobra.OnInitialize(config.InitConfig)

	rootCmd.PersistentFlags().StringVarP(&config.ConfigFile, "config", "c", "", "Config File")
}
