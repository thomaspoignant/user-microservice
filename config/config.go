package config

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/spf13/viper"
)

// LoadConfigFile load configuration from YAML file
func LoadConfigFile() {
	if os.Getenv("ENV") == "" && os.Getenv("TEST") != "true" {
		setLocalConfig()
	}
	viper.AutomaticEnv()
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("GIN_MODE", gin.ReleaseMode)
	viper.SetDefault("RUNNING_MODE", "api")
	viper.SetDefault("AWS_REGION", "eu-west-1")
	viper.SetDefault("DYNAMODB_ENDPOINT", "dynamodb.eu-west-1.amazonaws.com")
	viper.SetDefault("DYNAMODB_TABLE_NAME", "user")
}

// setLocalConfig surcharge configuration to execute app locally
func setLocalConfig() {
	os.Setenv("GIN_MODE", gin.DebugMode)
	os.Setenv("DYNAMODB_ENDPOINT", "http://localhost:9000")
}
