package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// InitConfig initializes the application configuration
func InitConfig(logger *otelzap.Logger) Config {
	viper.SetConfigName(".env") // name of config file (without extension)
	viper.SetConfigType("env")  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")    // optionally look for config in the working directory

	// Set up environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	// Set defaults
	viper.SetDefault("db_url", "")
	viper.SetDefault("redis_url", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No .env file found, using environment variables")
		} else {
			fmt.Printf("Error reading config file: %v\n", err)
		}
	}

	// Debug logging
	fmt.Printf("Loading configuration...\n")
	fmt.Printf("DB_URL: %s\n", viper.GetString("db_url"))
	fmt.Printf("REDIS_URL: %s\n", viper.GetString("redis_url"))
	fmt.Printf("FOOTBALL_API_KEY length: %d\n", len(viper.GetString("football_api_key")))

	return Config{
		DB_URL:           viper.GetString("db_url"),
		FOOTBALL_API_KEY: viper.GetString("football_api_key"),
		REDIS_URL:        viper.GetString("redis_url"),
	}
}
