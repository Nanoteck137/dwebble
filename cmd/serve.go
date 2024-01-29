package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Serve\n") 
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
