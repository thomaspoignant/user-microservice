package config

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ProjectPath is the home of the project
var ProjectPath = os.Getenv("GOPATH") + "/src/github.com/thomaspoignant/user-microservice"

// LoadConfigFile load configuration from YAML file
func LoadConfigFile() {
	if os.Getenv("ENV") == "" {
		configFileName := composeConfigFileName()
		log.Infof("Trying to load file : %s", configFileName)
		viper.SetConfigName(configFileName)
		viper.AddConfigPath(ProjectPath + "/config/")
		if err := viper.ReadInConfig(); err != nil {
			panic("Error reading config file : " + err.Error())
		}
	} else {
		log.Info("Loading config from environnement variables")
		viper.AutomaticEnv()
	}
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("GIN_MODE", gin.ReleaseMode)
	viper.SetDefault("RUNNING_MODE", "api")
	viper.SetDefault("AWS_REGION", "eu-west-1")
	viper.SetDefault("DYNAMODB_ENDPOINT", "dynamodb.eu-west-1.amazonaws.com")
}

// determine the name of the config file
func composeConfigFileName() string {
	env := os.Getenv("ENV")
	test := os.Getenv("TEST")
	var configFileName []string
	configFileName = append(configFileName, "config")

	if strings.Compare("true", test) == 0 {
		configFileName = append(configFileName, "test")
	}

	if strings.Compare("", env) == 0 {
		configFileName = append(configFileName, "default")
	} else {
		configFileName = append(configFileName, env)
	}

	fileName := strings.Join(configFileName, ".")
	return fileName
}
