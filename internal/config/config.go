package config

import (
	"context"
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
	viper.SetDefault("db_url", "")
	viper.SetDefault("port", "8080")

	// automatically load matching envs
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	viper.AutomaticEnv()

	// the following envs are not automatic because they didn't match the key structure
	// _ = viper.BindEnv("http.cookie_hashkey", "COOKIE_HASHKEY")
	// _ = viper.BindEnv("http.port", "PORT")
	// _ = viper.BindEnv("http.secure_cookie", "COOKIE_SECURE")
	// _ = viper.BindEnv("http.backend_cookie_name", "SECURE_COOKIE_NAME")
	// _ = viper.BindEnv("http.session_cookie_name", "SESSION_COOKIE_NAME")
	// _ = viper.BindEnv("http.frontend_cookie_name", "FRONTEND_COOKIE_NAME")
	// _ = viper.BindEnv("http.domain", "APP_DOMAIN")
	// _ = viper.BindEnv("http.path_prefix", "PATH_PREFIX")
	// _ = viper.BindEnv("config.allowedPointValues", "CONFIG_POINTS_ALLOWED")
	// _ = viper.BindEnv("config.defaultPointValues", "CONFIG_POINTS_DEFAULT")
	// _ = viper.BindEnv("config.show_warrior_rank", "CONFIG_SHOW_RANK")
	// _ = viper.BindEnv("auth.header.usernameHeader", "AUTH_HEADER_USERNAME_HEADER")
	// _ = viper.BindEnv("auth.header.emailHeader", "AUTH_HEADER_EMAIL_HEADER")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			logger.Ctx(context.Background()).Fatal(err.Error())
		} else {
			// Config file was found but another error was produced
			logger.Ctx(context.Background()).Fatal(err.Error())
		}
	}

	var c Config
	viper.Unmarshal(&c)
	// if err != nil {
	// 	logger.Ctx(context.Background()).Fatal(err.Error())
	// }
	// Config file found and successfully parsed

	return c
}
