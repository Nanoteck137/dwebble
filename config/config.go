package config

import (
	"os"

	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/types"
	"github.com/spf13/viper"
)

type Config struct {
	ListenAddr      string `mapstructure:"listen_addr"`
	DataDir         string `mapstructure:"data_dir"`
	LibraryDir      string `mapstructure:"library_dir"`
	Username        string `mapstructure:"username"`
	InitialPassword string `mapstructure:"initial_password"`
	JwtSecret       string `mapstructure:"jwt_secret"`
}

func (c *Config) WorkDir() types.WorkDir {
	return types.WorkDir(c.DataDir)
}

func setDefaults() {
	viper.SetDefault("listen_addr", ":3000")
	viper.BindEnv("data_dir")
	viper.BindEnv("library_dir")
	viper.BindEnv("username")
	viper.BindEnv("initial_password")
	viper.BindEnv("jwt_secret")
}

func validateConfig(config *Config) {
	hasError := false

	validate := func(expr bool, msg string) {
		if expr {
			log.Error("Config Validation", "err", msg)
			hasError = true
		}
	}

	// NOTE(patrik): Has default value, here for completeness
	validate(config.ListenAddr == "", "listen_addr needs to be set")
	validate(config.DataDir == "", "data_dir needs to be set")
	validate(config.LibraryDir == "", "library_dir needs to be set")
	validate(config.Username == "", "username needs to be set")
	validate(config.InitialPassword == "", "initial_password needs to be set")
	validate(config.JwtSecret == "", "jwt_secret needs to be set")

	if hasError {
		log.Fatal("Config not valid")
	}
}

var ConfigFile string
var LoadedConfig Config

func InitConfig() {
	setDefaults()

	if ConfigFile != "" {
		viper.SetConfigFile(ConfigFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix(dwebble.AppName)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Warn("Failed to load config", "err", err)
	}

	err = viper.Unmarshal(&LoadedConfig)
	if err != nil {
		log.Error("Failed to unmarshal config: ", err)
		os.Exit(-1)
	}


	hide := func(s string) string {
		var res string
		for i := 0; i < len(s); i++ {
			res += "*"
		}

		return res
	}

	configCopy := LoadedConfig
	configCopy.JwtSecret = hide(configCopy.JwtSecret)
	configCopy.InitialPassword = hide(configCopy.InitialPassword)

	log.Debug("Current Config", "config", configCopy)

	validateConfig(&LoadedConfig)
}
