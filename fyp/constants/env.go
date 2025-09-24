package constants

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	PowNumberRule string `mapstructure:"POW_NUMBER_RULE"`
}

func NewEnv() *Env {
	viper.SetConfigFile("../.env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}
	viper.AutomaticEnv()

	var env Env
	if err := viper.Unmarshal(&env); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	return &env
}

func (e *Env) GetValueForKey(key string) string {
	// Use viper directly to get the value of the key
	return viper.GetString(key)
}
