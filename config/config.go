package config

import (
	"fund-o/api-server/internal/datasource"
	"github.com/spf13/viper"
)

type AppConfig struct {
	AppEnv   string `mapstructure:"APP_ENV"`
	LogLevel string `mapstructure:"LOG_LEVEL"`
	GinMode  string `mapstructure:"GIN_MODE"`
	ApiServerConfig
	datasource.DatasourceConfig
}

func LoadAppConfig(path string) (config AppConfig, err error) {
	makeDefaultAppConfig()
	//viper.AddConfigPath(path)
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	//viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func makeDefaultAppConfig() {
	// Set default values for app configuration
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("LOG_REQUEST", true)
	viper.SetDefault("GIN_MODE", "debug")

	// Set default values for api server configuration
	viper.SetDefault("ApiServerConfig.APP_HOST", "localhost")
	viper.SetDefault("ApiServerConfig.APP_PORT", "3000")
	viper.SetDefault("ApiServerConfig.APP_PATH_PREFIX", "/api/v1")
	viper.SetDefault("ApiServerConfig.APP_REQUEST_ID_HEADER", "X-Request-Id")
	viper.SetDefault("ApiServerConfig.APP_TRUST_PROXY", "10.0.0.0/8,172.16.0.0/12,192.168.0.0/16,fd00::/8")
	viper.SetDefault("ApiServerConfig.APP_CORS_ENABLED", true)
	viper.SetDefault("ApiServerConfig.APP_CORS_ALLOWED_ORIGIN", "*")
	viper.SetDefault("ApiServerConfig.APP_CORS_ALLOWED_CREDENTIALS", true)
	viper.SetDefault("ApiServerConfig.APP_CORS_MAX_AGE", 300)
	viper.SetDefault("ApiServerConfig.APP_READ_ONLY", false)
	viper.SetDefault("ApiServerConfig.LOG_REQUEST", true)

	// Set default values for sql db configuration
	viper.SetDefault("DatasourceConfig.SqlDBConfig.SQL_HOST", "localhost")
	viper.SetDefault("DatasourceConfig.SqlDBConfig.SQL_USERNAME", "docker")
	viper.SetDefault("DatasourceConfig.SqlDBConfig.SQL_PASSWORD", "secret")
	viper.SetDefault("DatasourceConfig.SqlDBConfig.SQL_PORT", 5432)
	viper.SetDefault("DatasourceConfig.SqlDBConfig.SQL_DATABASE", "fundo")
}
