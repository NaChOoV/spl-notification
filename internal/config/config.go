package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type EnvironmentConfig struct {
	Port                   string `env:"PORT,default=4002"`
	AuthString             string `env:"AUTH_STRING,required"`
	TursoBaseUrl           string `env:"TURSO_DATABASE_URL,required"`
	TursoAuthToken         string `env:"TURSO_AUTH_TOKEN,required"`
	NotificationBaseUrl    string `env:"NOTIFICATION_BASE_URL,required"`
	NotificationUsername   string `env:"NOTIFICATION_USERNAME,required"`
	NotificationPassword   string `env:"NOTIFICATION_PASSWORD,required"`
	AccessServiceBaseUrl   string `env:"ACCESS_SERVICE_BASE_URL,required"`
	AccessServiceAuthToken string `env:"ACCESS_SERVICE_AUTH_TOKEN,required"`
	DebugMode              bool   `env:"DEBUG_MODE"`
	Environment            string `env:"ENVIRONMENT,default=LOCAL"`

	Zone string `env:"ZONE"`

	// Source Service
	SourceBaseUrl    string `env:"SOURCE_BASE_URL,required"`
	SourceAuthString string `env:"SOURCE_AUTH_STRING,required"`

	// Google Cloud Pub/Sub
	PubSubProjectID      string `env:"PUBSUB_PROJECT_ID,required"`
	PubSubTopicID        string `env:"PUBSUB_TOPIC_ID,required"`
	PubSubSubscriptionID string `env:"PUBSUB_SUBSCRIPTION_ID,required"`
}

var envConfig *EnvironmentConfig

func NewEnviromentConfig() *EnvironmentConfig {
	envConfig = &EnvironmentConfig{}

	var err error
	enviroment := os.Getenv("ENVIRONMENT")
	if enviroment == "LOCAL" || enviroment == "" {
		err = godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
			panic(err)
		}
	}

	// Port
	envConfig.Port = os.Getenv("PORT")
	if envConfig.Port == "" {
		envConfig.Port = "4001"
	}

	// AuthString
	envConfig.AuthString = os.Getenv("AUTH_STRING")

	// Turso
	envConfig.TursoBaseUrl = os.Getenv("TURSO_DATABASE_URL")
	envConfig.TursoAuthToken = os.Getenv("TURSO_AUTH_TOKEN")

	// Notification Service
	envConfig.NotificationBaseUrl = os.Getenv("NOTIFICATION_BASE_URL")
	envConfig.NotificationUsername = os.Getenv("NOTIFICATION_USERNAME")
	envConfig.NotificationPassword = os.Getenv("NOTIFICATION_PASSWORD")

	// DebugMode
	if os.Getenv("DEBUG_MODE") == "true" {
		envConfig.DebugMode = true
	} else {
		envConfig.DebugMode = false
	}

	// Zone
	envConfig.Zone = os.Getenv("ZONE")
	if envConfig.Zone == "" {
		envConfig.Zone = "GMT-3"
	}

	// Source Service
	envConfig.SourceBaseUrl = os.Getenv("SOURCE_BASE_URL")
	envConfig.SourceAuthString = os.Getenv("SOURCE_AUTH_STRING")

	// Access Service
	envConfig.AccessServiceBaseUrl = os.Getenv("ACCESS_SERVICE_BASE_URL")
	envConfig.AccessServiceAuthToken = os.Getenv("ACCESS_SERVICE_AUTH_TOKEN")

	// Google Cloud Pub/Sub
	envConfig.PubSubProjectID = os.Getenv("PUBSUB_PROJECT_ID")
	envConfig.PubSubTopicID = os.Getenv("PUBSUB_TOPIC_ID")
	envConfig.PubSubSubscriptionID = os.Getenv("PUBSUB_SUBSCRIPTION_ID")

	printEnvironmentConfig(*envConfig)
	return envConfig
}

func printEnvironmentConfig(config EnvironmentConfig) {
	v := reflect.ValueOf(config)
	typeOfConfig := v.Type()

	fmt.Println("Environments:")
	for i := 0; i < v.NumField(); i++ {
		fmt.Printf("  %s: %v\n", typeOfConfig.Field(i).Name, v.Field(i).Interface())
	}
}
