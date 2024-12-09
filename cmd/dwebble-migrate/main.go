package main

import (
	"strings"

	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/apis"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/spf13/cobra"
)

var AppName = dwebble.AppName + "-migrate"

var rootCmd = &cobra.Command{
	Use:     AppName + " <IN>",
	Version: dwebble.Version,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.SetVersionTemplate(dwebble.VersionTemplate(AppName))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Failed to run root command", "err", err)
	}
}

func main() {
	Execute()
}
