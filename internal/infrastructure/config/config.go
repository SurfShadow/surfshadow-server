package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `mapstructure:"SERVER" validate:"required"`
	DB      DBConfig      `mapstructure:"DB" validate:"required"`
	Logger  LoggerConfig  `mapstructure:"LOGGER" validate:"required"`
	Metrics MetricsConfig `mapstructure:"METRICS" validate:"required"`
	AuthAPI AuthAPIConfig `mapstructure:"API" validate:"required"`
}

type ServerConfig struct {
	AppVersion string `mapstructure:"APP_VERSION" validate:"required"`
	Port       uint16 `mapstructure:"PORT" validate:"required,gt=0"`
	SSL        bool   `mapstructure:"SSL"`
}

type DBConfig struct {
	Host         string `mapstructure:"POSTGRES_HOST" validate:"required"`
	Port         uint16 `mapstructure:"POSTGRES_PORT" validate:"required,gt=0"`
	User         string `mapstructure:"POSTGRES_USER" validate:"required"`
	Password     string `mapstructure:"POSTGRES_PASSWORD" validate:"required"`
	DataBaseName string `mapstructure:"POSTGRES_DB" validate:"required"`
	PgDriver     string `mapstructure:"POSTGRES_DRIVER" validate:"required"`
}

type LoggerConfig struct {
	Level string `mapstructure:"LOG_LEVEL" validate:"required,oneof=debug info warn error dpanic panic fatal"`
}

type MetricsConfig struct {
	Port int16 `mapstructure:"PORT" validate:"required,gt=0"`
}

type AuthAPIConfig struct {
	APIKey string `mapstructure:"AUTH_KEY" validate:"required"`
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

	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}

	return &config, nil
}
