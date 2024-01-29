// +build dev

package cmd

import "github.com/spf13/cobra"

var mockCmd = &cobra.Command{
	Use: "mock",
	Short: "Generate mock data",
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(mockCmd)
}
