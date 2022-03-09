package config

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func SetupConfig() {
	// Reading .env file
	viper.SetConfigName(".env")
	viper.AddConfigPath("./")
	viper.SetConfigType("env")

	if err := viper.MergeInConfig(); err != nil {
		log.Warn().Err(err).Msg("couldn't find .env file to read variables from")
	}

	// Reading env variables
	viper.AutomaticEnv()
}

func SetupLogger() {
	if os.Getenv("ENV") == "local" {
		output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"}
		output.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s\t", i)
		}

		log.Logger = log.Output(output)
	}
}
