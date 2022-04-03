package config

import (
	"os"

	"github.com/skamranahmed/markeventas/pkg/log"
	"github.com/spf13/viper"
)

// AppEnvironment : string wrapper for environment name
type AppEnvironment string

func (e AppEnvironment) IsLocal() bool {
	return e == AppEnvironmentLocal
}

var (
	// Database Credentials
	DbHost     string
	DbUser     string
	DbName     string
	DbPassword string
	DbPort     string

	// Server
	ServerPort string

	// Token
	TokenSecretSigningKey string

	// Google API
	GoogleAppClientSecret string

	// Twitter Gcal Event Login App Credentials
	TwitterLoginAppApiKey       string
	TwitterLoginAppApiKeySecret string

	// Twitter Gcal Event Bot Credentials
	TwitterBotApiKey            string
	TwitterBotApiKeySecret      string
	TwitterBotAccessToken       string
	TwitterBotAccessTokenSecret string

	// Environment
	Environment AppEnvironment

	// slice of all app environments except the `local`` env
	AppEnvironemnts = []AppEnvironment{
		AppEnvironmentStaging,
		AppEnvironmentSandbox,
		AppEnvironmentProduction,
	}
)

const (
	// AppEnvironmentLocal : local env
	AppEnvironmentLocal = AppEnvironment("local")

	// AppEnvironmentLocal : staging env
	AppEnvironmentStaging = AppEnvironment("staging")

	// AppEnvironmentLocal : sandbox env
	AppEnvironmentSandbox = AppEnvironment("sandbox")

	// AppEnvironmentLocal : production env
	AppEnvironmentProduction = AppEnvironment("production")

	// ConfigFileName : localConfig.yaml
	ConfigFileName string = "localConfig"

	// ConfigFileType : yaml
	ConfigFileType string = "yaml"
)

func init() {
	SetConfigFromViper()
}

func SetConfigFromViper() {
	Environment = getCurrentHostEnvironment()
	log.Infof("🚀 Current Host Environment: %s\n", Environment)

	// if env is local, we set the env variables using the config file
	if Environment.IsLocal() {
		setEnvironmentVarsFromConfig()
	}

	// fetch the env vars and store in variables
	// Database Credentials
	DbHost = os.Getenv("DB_HOST")
	DbUser = os.Getenv("DB_USER")
	DbName = os.Getenv("DB_NAME")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbPort = os.Getenv("DB_PORT")

	// Server
	ServerPort = os.Getenv("SERVER_PORT")

	// Token
	TokenSecretSigningKey = os.Getenv("TOKEN_SECRET_SIGNING_KEY")

	// Google API
	GoogleAppClientSecret = os.Getenv("GOOGLE_APP_CLIENT_SECRET")

	// Twitter Gcal Event Login App Credentials
	TwitterLoginAppApiKey = os.Getenv("TWITTER_LOGIN_APP_API_KEY")
	TwitterLoginAppApiKeySecret = os.Getenv("TWITTER_LOGIN_APP_API_KEY_SECRET")

	// Twitter Gcal Event Bot Credentials
	TwitterBotApiKey = os.Getenv("TWITTER_BOT_API_KEY")
	TwitterBotApiKeySecret = os.Getenv("TWITTER_BOT_API_KEY_SECRET")
	TwitterBotAccessToken = os.Getenv("TWITTER_BOT_ACCESS_TOKEN")
	TwitterBotAccessTokenSecret = os.Getenv("TWITTER_BOT_ACCESS_TOKEN_SECRET")
}

func setEnvironmentVarsFromConfig() {
	baseProjectPath, _ := os.Getwd()

	// add the path of the config file
	viper.AddConfigPath(baseProjectPath + "/config/")
	// set the config file name
	viper.SetConfigName(ConfigFileName)
	// set the config file type
	viper.SetConfigType(ConfigFileType)

	viper.AutomaticEnv()

	// read the env vars from the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("unable to read env vars from config file")
	}

	/*
		Step 1. Get the env vars from viper
		Step 2. Set the host OS env vars
	*/

	// Database Credentials
	dbHost := viper.GetString("DB_HOST")
	dbUser := viper.GetString("DB_USER")
	dbName := viper.GetString("DB_NAME")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbPort := viper.GetString("DB_PORT")
	os.Setenv("DB_HOST", dbHost)
	os.Setenv("DB_USER", dbUser)
	os.Setenv("DB_NAME", dbName)
	os.Setenv("DB_PASSWORD", dbPassword)
	os.Setenv("DB_PORT", dbPort)

	// Server
	serverPort := viper.GetString("SERVER_PORT")
	os.Setenv("SERVER_PORT", serverPort)

	// Token
	tokenSecretSigningKey := viper.GetString("TOKEN_SECRET_SIGNING_KEY")
	os.Setenv("TOKEN_SECRET_SIGNING_KEY", tokenSecretSigningKey)

	// Google API
	googleAppClientSecret := viper.GetString("GOOGLE_APP_CLIENT_SECRET")
	os.Setenv("GOOGLE_APP_CLIENT_SECRET", googleAppClientSecret)

	// Twitter Gcal Event Login App Credentials
	twitterLoginAppApiKey := viper.GetString("TWITTER_LOGIN_APP_API_KEY")
	twitterLoginAppApiKeySecret := viper.GetString("TWITTER_LOGIN_APP_API_KEY_SECRET")
	os.Setenv("TWITTER_LOGIN_APP_API_KEY", twitterLoginAppApiKey)
	os.Setenv("TWITTER_LOGIN_APP_API_KEY_SECRET", twitterLoginAppApiKeySecret)

	// Twitter Gcal Event Bot Credentials
	twitterBotApiKey := viper.GetString("TWITTER_BOT_API_KEY")
	twitterBotApiKeySecret := viper.GetString("TWITTER_BOT_API_KEY_SECRET")
	twitterBotAccessToken := viper.GetString("TWITTER_BOT_ACCESS_TOKEN")
	twitterBotAccessTokenSecret := viper.GetString("TWITTER_BOT_ACCESS_TOKEN_SECRET")
	os.Setenv("TWITTER_BOT_API_KEY", twitterBotApiKey)
	os.Setenv("TWITTER_BOT_API_KEY_SECRET", twitterBotApiKeySecret)
	os.Setenv("TWITTER_BOT_ACCESS_TOKEN", twitterBotAccessToken)
	os.Setenv("TWITTER_BOT_ACCESS_TOKEN_SECRET", twitterBotAccessTokenSecret)

}

func getCurrentHostEnvironment() AppEnvironment {
	currentHostEnvironment := os.Getenv("ENVIRONMENT")
	for _, env := range AppEnvironemnts {
		if env == AppEnvironment(currentHostEnvironment) {
			return env
		}
	}
	// if env not found return `local`` env
	return AppEnvironmentLocal
}
