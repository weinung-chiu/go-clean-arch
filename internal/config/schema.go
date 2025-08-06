package config

type AppConfig struct {
	Env       string `mapstructure:"APP_ENV"`
	LogLevel  string `mapstructure:"APP_LOG_LEVEL"`
	ApiPort   int    `mapstructure:"API_PORT"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
	JWTExpiry int    `mapstructure:"JWT_EXPIRY_MINUTES"` // Token expiry in minutes
	DatabaseDSN string `mapstructure:"DATABASE_DSN"`
}
