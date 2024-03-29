package config

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"time"
)

type EnvConfigs struct {
	ServerAddressPort       int           `mapstructure:"SERVER_ADDRESS_PORT"`
	ServerTimeoutRead       time.Duration `mapstructure:"SERVER_TIMEOUT_READ"`
	ServerTimoutWrite       time.Duration `mapstructure:"SERVER_TIMEOUT_WRITE"`
	ServerTimeoutIdle       time.Duration `mapstructure:"SERVER_TIMEOUT_IDLE"`
	DBDsn                   string        `mapstructure:"DB_DSN"`
	DBMaxConnectionLifeTime time.Duration `mapstructure:"DB_MAX_CONNECTION_LIFETIME"`
	DBMaxConnection         int32         `mapstructure:"DB_MAX_CONNECTION"`
	DBMinConnection         int32         `mapstructure:"DB_MIN_CONNECTION"`
	DBMaxConnectionIdleTime time.Duration `mapstructure:"DB_MAX_CONNECTION_IDLE_TIME"`

	AccessTokenPrivateKey string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey  string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiredIn  time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`

	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	RefreshTokenExpiredIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`

	RedisAddressAndPort string `mapstructure:"REDIS_ADDRESS_AND_PORT"`
	RedisPassword       string `mapstructure:"REDIS_PASSWORD"`
	RedisDB             int    `mapstructure:"REDIS_DB"`
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
