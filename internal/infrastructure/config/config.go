package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `mapstructure:"SERVER"`
	DB      DBConfig      `mapstructure:"DB"`
	Logger  LoggerConfig  `mapstructure:"LOGGER"`
	Metrics MetricsConfig `mapstructure:"METRICS"`
	AuthAPI AuthAPIConfig `mapstructure:"API"`
}

type ServerConfig struct {
	AppVersion string `mapstructure:"APP_VERSION"`
	Port       uint16 `mapstructure:"PORT"`
	SSL        bool   `mapstructure:"SSL"`
}

type DBConfig struct {
	Host         string `mapstructure:"POSTGRES_HOST"`
	Port         uint16 `mapstructure:"POSTGRES_PORT"`
	User         string `mapstructure:"POSTGRES_USER"`
	Password     string `mapstructure:"POSTGRES_PASSWORD"`
	DataBaseName string `mapstructure:"POSTGRES_DB"`
	PgDriver     string `mapstructure:"POSTGRES_DRIVER"`
}

type LoggerConfig struct {
	Level string `mapstructure:"LOG_LEVEL"`
}

type MetricsConfig struct {
	Port int16 `mapstructure:"PORT"`
}

type AuthAPIConfig struct {
	APIKey string `mapstructure:"AUTH_KEY"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.SetDefault("SERVER.SERVER_PORT", 8080)
	viper.SetDefault("SERVER.SERVER_SSL", false)

	viper.SetDefault("DB.POSTGRES_HOST", "localhost")
	viper.SetDefault("DB.POSTGRES_PORT", 5432)
	viper.SetDefault("DB.POSTGRES_USER", "admin")
	viper.SetDefault("DB.POSTGRES_PASSWORD", "admin")
	viper.SetDefault("DB.POSTGRES_DB", "database")
	viper.SetDefault("DB.POSTGRES_DRIVER", "pgx")

	viper.SetDefault("LOGGER.LOG_LEVEL", "info")

	viper.SetDefault("METRICS_PORT", 9090)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &config, nil
}
