package api

import (
	"github.com/spf13/viper"
)

// Init the configuration of the APIs
func Init() {
	router := SetupRouter()
	router.Run(":" + viper.GetString("APP_PORT"))
}
