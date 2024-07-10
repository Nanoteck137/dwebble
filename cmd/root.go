package cmd

import (
	"fmt"

	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     config.AppName,
	Version: config.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Failed to run root command", "err", err)
	}
}

func versionTemplate() string {
	return fmt.Sprintf(
		"%s: %s (%s)\n",
		config.AppName, config.Version, config.Commit)
}

func init() {
	rootCmd.SetVersionTemplate(versionTemplate())

	cobra.OnInitialize(config.InitConfig)

	rootCmd.PersistentFlags().StringVarP(&config.ConfigFile, "config", "c", "", "Config File")
}
