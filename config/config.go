package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort    string
	DBSupabase string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load() *Config {
	// ENV support
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Load .env if exists (local only)
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	viper.SetDefault("APP_PORT", "8080")

	return &Config{
		AppPort:    viper.GetString("APP_PORT"),
		DBSupabase: viper.GetString("DB_SUPABASE_MAIN"),
	}
}
