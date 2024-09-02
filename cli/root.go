package cli

import (
	"fmt"

	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: dwebble.AppName,
	Version: dwebble.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Failed to run root command", "err", err)
	}
}

func versionTemplate() string {
	return fmt.Sprintf(
		"%s: %s (%s)\n",
		dwebble.AppName, dwebble.Version, dwebble.Commit)
}

func init() {
	rootCmd.SetVersionTemplate(versionTemplate())

	cobra.OnInitialize(config.InitConfig)

	rootCmd.PersistentFlags().StringVarP(&config.ConfigFile, "config", "c", "", "Config File")
}
