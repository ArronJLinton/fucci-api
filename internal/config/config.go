package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// InitConfig initializes the application configuration
func InitConfig(logger *otelzap.Logger) Config {
	viper.SetConfigName(".env") // name of config file (without extension)
	viper.SetConfigType("env")  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	// viper.SetConfigFile("config")

	// Set up environment variables first
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	// Explicitly bind PORT environment variable
	viper.BindEnv("port", "PORT")

	// Set defaults only if env vars are not present
	if os.Getenv("PORT") != "" {
		viper.SetDefault("port", os.Getenv("PORT"))
	} else {
		viper.SetDefault("port", "8080")
	}
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
	// if err != nil {
	// 	logger.Ctx(context.Background()).Fatal(err.Error())
	// }
	// Config file found and successfully parsed

	// Debug log all environment variables
	logger.Ctx(context.Background()).Info(fmt.Sprintf("Raw PORT env: %s", os.Getenv("PORT")))
	logger.Ctx(context.Background()).Info(fmt.Sprintf("Viper PORT: %s", viper.GetString("port")))
	logger.Ctx(context.Background()).Info(fmt.Sprintf("Config PORT: %s", c.PORT))
	logger.Ctx(context.Background()).Info(fmt.Sprintf("DB_URL: %s", c.DB_URL))
	logger.Ctx(context.Background()).Info(fmt.Sprintf("REDIS_URL: %s", c.REDIS_URL))

	return c
}
