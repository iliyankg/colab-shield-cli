package config

import (
	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config from working directory. Make sure config.json is present.")
	}

	// TODO: Validate config.
}

func ServerHost() string {
	return viper.GetString("server.host")
}

func ServerPortGrpc() int {
	return viper.GetInt("server.ports.grpc")
}

func ProjectId() string {
	return viper.GetString("project.id")
}

func Extensions() []string {
	return viper.GetStringSlice("extensions")
}

func IgnorePaths() []string {
	return viper.GetStringSlice("ignore")
}
