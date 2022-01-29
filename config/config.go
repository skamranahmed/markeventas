package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	// Twitter Gcal Event Login App Credentials
	TwitterGcalEventLoginAppApiKey       string `mapstructure:"TWITTER_CREATE_GCAL_EVENT_LOGIN_APP_API_KEY"`
	TwitterGcalEventLoginAppApiKeySecret string `mapstructure:"TWITTER_CREATE_GCAL_EVENT_LOGIN_APP_API_KEY_SECRET"`

	// Twitter Gcal Event Bot Credentials
	TwitterGcalEventBotApiKey       string `mapstructure:"TWITTER_CREATE_GCAL_EVENT_BOT_API_KEY"`
	TwitterGcalEventBotApiKeySecret string `mapstructure:"TWITTER_CREATE_GCAL_EVENT_BOT_API_KEY_SECRET"`

	// Database Credentials
	DbHost     string `mapstructure:"DB_HOST"`
	DbUser     string `mapstructure:"DB_USER"`
	DbName     string `mapstructure:"DB_NAME"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbPort     string `mapstructure:"DB_PORT"`

	TokenSecretSigningKey string `mapstructure:"TOKEN_SECRET_SIGNING_KEY"`

	// Server
	ServerPort string `mapstructure:"SERVER_PORT"`
}

func LoadConfig() (*Config, error) {
	// provide the config file name
	viper.SetConfigName("localConfig")

	// provide the config file path
	configFilePath, err := getConfigFilePath()
	if err != nil {
		log.Fatalf("Error in getting the config file path: %v\n", err)
	}
	viper.AddConfigPath(configFilePath)

	// Find and read the config file
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error in reading config file: %v\n", err)
	}

	// find matching env vars and load them in Viper
	// this will override the values for that specific env var mentioned in the config file
	viper.AutomaticEnv()

	var conf Config
	err = viper.Unmarshal(&conf)
	return &conf, err
}

func getConfigFilePath() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		log.Printf("Error in getting the workind directory path: %v\n", err)
		return "", err
	}
	configFilePath := fmt.Sprintf("%s/config/", path)
	return configFilePath, nil
}
