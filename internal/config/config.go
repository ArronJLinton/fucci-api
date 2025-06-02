package config

import (
	"context"
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
			// Config file not found, continue with env vars
			logger.Ctx(context.Background()).Info("Config file .env not found, using environment variables")
		} else {
			// Config file was found but another error was produced
			logger.Ctx(context.Background()).Fatal(err.Error())
		}
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		logger.Ctx(context.Background()).Fatal(fmt.Sprintf("Failed to unmarshal config: %v", err))
	}

	// Debug log configuration
	logger.Ctx(context.Background()).Info(fmt.Sprintf("DB_URL: %s", c.DB_URL))
	logger.Ctx(context.Background()).Info(fmt.Sprintf("REDIS_URL: %s", c.REDIS_URL))

	return c
}
