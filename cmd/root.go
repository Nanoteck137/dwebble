package cmd

import (
	"fmt"
	"log"

	"github.com/nanoteck137/dwebble/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     config.AppName,
	Version: config.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
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
