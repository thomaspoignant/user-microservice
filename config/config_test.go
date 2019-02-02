package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// TestConfigLocal is testing that we read the value from the config file (config.default.yaml)
func TestConfigLocal(t *testing.T) {
	cleanEnvVarForLoadingConfigFile()
	expected := "debug"
	LoadConfigFile()
	got := viper.GetString("GIN_MODE")
	assert.Equal(t, expected, got)
}

// TestConfigRelease is testing we reading config from environnement variable in release
func TestConfigReleaseFromEnvVar(t *testing.T) {
	cleanEnvVarForLoadingConfigFile()
	//setting the environnement variable "ENV" to DEV
	os.Setenv("ENV", "DEV")
	envVarName := "APP_PORT"
	expected := "8585"
	os.Setenv(envVarName, expected)
	LoadConfigFile()
	actual := viper.GetString(envVarName)
	assert.Equal(t, expected, actual)
	//clean env variables
	os.Setenv("APP_PORT", "")
}

// TestConfigReleaseNoValueUsingDefault is testing we use default value if no env variable are set
func TestConfigReleaseNoValueUsingDefault(t *testing.T) {
	cleanEnvVarForLoadingConfigFile()
	//setting the environnement variable "ENV" to DEV
	os.Setenv("ENV", "DEV")
	envVarName := "APP_PORT"
	expected := "8080"
	LoadConfigFile()
	got := viper.GetString(envVarName)
	assert.Equal(t, expected, got)
}

func cleanEnvVarForLoadingConfigFile() {
	os.Setenv("TEST", "")
	os.Setenv("ENV", "")
}
