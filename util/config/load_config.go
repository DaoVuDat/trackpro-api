package config

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"time"
)

type EnvConfigs struct {
	ServerAddressPort int           `mapstructure:"SERVER_ADDRESS_PORT"`
	ServerTimeoutRead time.Duration `mapstructure:"SERVER_TIMEOUT_READ"`
	ServerTimoutWrite time.Duration `mapstructure:"SERVER_TIMEOUT_WRITE"`
	ServerTimeoutIdle time.Duration `mapstructure:"SERVER_TIMEOUT_IDLE"`
	DBDsn             string        `mapstructure:"DB_DSN"`
}

func LoadEnvConfigs(logger *zerolog.Logger, path string) *EnvConfigs {
	var envConfig *EnvConfigs
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Panic().Err(err).Msg("Error reading env file")
		panic(1)
	}

	if err := viper.Unmarshal(&envConfig); err != nil {
		logger.Panic().Err(err).Msg("Cannot unmarshal env config")
		panic(1)
	}

	return envConfig
}
