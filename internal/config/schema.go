package config

type AppConfig struct {
	Env      string `mapstructure:"APP_ENV"`
	LogLevel string `mapstructure:"APP_LOG_LEVEL"`
	ApiPort  int    `mapstructure:"API_PORT"`
}
