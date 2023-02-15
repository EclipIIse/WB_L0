package config

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string `mapstructure:"HOST"`
}

type DbConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Database string `mapstructure:"DB_NAME"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
}

type NATS struct {
	ClusterName string `mapstructure:"NATS_CLUSTER_NAME"`
	ClientID    string `mapstructure:"NATS_CLIENT_ID"`
	SubjectName string `mapstructure:"NATS_SUBJECT_NAME"`
}

type Redis struct {
	Address    string `mapstructure:"REDIS_ADDRESS"`
	MaxRetries int    `mapstructure:"REDIS_RETRIES"`
}

type FrontConfig struct {
	Host string `mapstructure:"FRONT_HOST"`
}

func SetupSource() {
	viper.SetConfigType("env")
	viper.AutomaticEnv()
}

type Config struct {
	ServerConfig ServerConfig
	DbConfig     DbConfig
	NATS         NATS
	Redis        Redis
	FrontConfig  FrontConfig
}

func NewConfig() *Config {
	SetupSource()

	return &Config{
		ServerConfig: ServerConfig{viper.GetString("HOST")},
		DbConfig: DbConfig{
			viper.GetString("DB_HOST"),
			viper.GetString("DB_NAME"),
			viper.GetString("DB_USER"),
			viper.GetString("DB_PASSWORD"),
		},
		NATS: NATS{
			viper.GetString("NATS_CLUSTER_NAME"),
			viper.GetString("NATS_CLIENT_ID"),
			viper.GetString("NATS_SUBJECT_NAME"),
		},
		Redis: Redis{
			viper.GetString("REDIS_ADDRESS"),
			viper.GetInt("REDIS_RETRIES"),
		},
	}
}

func Logger() (logger zerolog.Logger) {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	return zerolog.New(output).With().Timestamp().Logger()
}
