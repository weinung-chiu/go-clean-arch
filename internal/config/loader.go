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
	
	// Try to load from .env file
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")
	
	// Read config file if it exists, but don't fail if it doesn't
	_ = v.ReadInConfig()
	
	setDefaultValues(v)

	var cfg AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal AppConfig failed: %w", err)
	}

	// Validate required fields
	if cfg.DatabaseDSN == "" {
		return nil, fmt.Errorf("DATABASE_DSN is required but not set")
	}

	return &cfg, nil
}

func setDefaultValues(v *viper.Viper) {
	v.SetDefault("APP_ENV", DefaultEnv)
	v.SetDefault("APP_LOG_LEVEL", DefaultLogLevel)
	v.SetDefault("API_PORT", DefaultPort)
}
