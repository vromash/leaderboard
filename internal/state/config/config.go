package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func SetupConfig() {
	env := os.Getenv("ENV")

	// Ignore configs directory if it doesn't exist
	if _, err := os.Stat("configs"); !os.IsNotExist(err) {
		// Reading config file
		viper.SetConfigName(env)
		viper.AddConfigPath("./configs")
		viper.SetConfigType("yml")

		if err := viper.ReadInConfig(); err != nil {
			log.Panic().Err(err).Msg("failed to read config file")
		}
	}

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

func GetResultsPerPage() int64 {
	var opt = viper.GetString("results_per_page")
	resultsPerPage, _ := strconv.ParseInt(opt, 10, 64)

	if resultsPerPage == 0 {
		return 10
	}

	return resultsPerPage
}
