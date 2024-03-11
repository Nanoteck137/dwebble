package cmd

import (
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use: "dwebble",
	Short: "Custom music server",
	Version: "v0.3.0",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config File")
	rootCmd.PersistentFlags().StringP("data-dir", "d", "", "Data Dir")
	viper.BindPFlag("data_dir", rootCmd.PersistentFlags().Lookup("data-dir"))
}

func setDefaults() {
	viper.SetDefault("listen_addr", ":3000")

	stateHome := os.Getenv("XDG_STATE_HOME")
	if stateHome == "" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		stateHome = path.Join(userHome, ".local", "state")
	} 

	dataDir := path.Join(stateHome, "dwebble")
	viper.SetDefault("data_dir", dataDir)
}

func initConfig() {
	setDefaults()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Failed to load config: ", err)
	}
}
