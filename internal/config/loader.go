package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

const (
	DefaultEnv      = "prod"
	DefaultLogLevel = "warn"
	DefaultPort     = 8080
)

func Load() (*AppConfig, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	setDefaultValues(v)

	var cfg AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal AppConfig failed: %w", err)
	}

	return &cfg, nil
}

func setDefaultValues(v *viper.Viper) {
	v.SetDefault("APP_ENV", DefaultEnv)
	v.SetDefault("APP_LOG_LEVEL", DefaultLogLevel)
	v.SetDefault("API_PORT", DefaultPort)
}
