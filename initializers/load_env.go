package initializers

import (
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/viper"
)

type Environment string

const (
	DevelopmentEnv Environment = "development"
	ProductionEnv  Environment = "production"
)

type Config struct {
	PORT           string      `mapstructure:"PORT"`
	DEV_URL        string      `mapstructure:"DEV_URL"`
	FRONTEND_URL   string      `mapstructure:"FRONTEND_URL"`
	HACKATHONS_URL string      `mapstructure:"HACKATHONS_URL"`
	ENV            Environment `mapstructure:"ENV"`
	JWT_SECRET     string      `mapstructure:"JWT_SECRET"`
	SOCKETS_URL    string      `mapstructure:"SOCKETS_URL"`
	LOGGER_URL     string      `mapstructure:"LOGGER_URL"`
	LOGGER_SECRET  string      `mapstructure:"LOGGER_SECRET"`
	LOGGER_TOKEN   string      `mapstructure:"LOGGER_TOKEN"`
}

var CONFIG Config

func LoadEnv() {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&CONFIG)
	if err != nil {
		log.Fatal(err)
	}

	requiredKeys := getRequiredKeys(CONFIG)
	missingKeys := checkMissingKeys(requiredKeys, CONFIG)

	if len(missingKeys) > 0 {
		err := fmt.Errorf("following environment variables not found: %v", missingKeys)
		log.Fatal(err)
	}

	if CONFIG.ENV != DevelopmentEnv && CONFIG.ENV != ProductionEnv {
		err := fmt.Errorf("invalid ENV value: %s", CONFIG.ENV)
		log.Fatal(err)
	}
}

func getRequiredKeys(config Config) []string {
	requiredKeys := []string{}
	configType := reflect.TypeOf(config)

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag != "" {
			requiredKeys = append(requiredKeys, tag)
		}
	}

	return requiredKeys
}

func checkMissingKeys(requiredKeys []string, config Config) []string {
	missingKeys := []string{}

	configValue := reflect.ValueOf(config)
	for _, key := range requiredKeys {
		value := configValue.FieldByName(key).String()
		if value == "" {
			missingKeys = append(missingKeys, key)
		}
	}

	return missingKeys
}
