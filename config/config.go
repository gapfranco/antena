package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config representa os parâmetros de configuração
type Config struct {
	TursoURL   string `mapstructure:"TURSO_URL"`
	TursoToken string `mapstructure:"TURSO_TOKEN"`
	Port       string `mapstructure:"ADDR"`
}

func GetConfig() (Config, error) {
	var config Config

	viper.SetConfigName("antena.conf")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("parsing config data: %w", err)
	}

	return config, nil
}
