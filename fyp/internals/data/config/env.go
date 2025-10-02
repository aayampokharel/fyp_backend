package config

import (
	err "project/package/errors"

	"github.com/spf13/viper"
)

type Env struct {
	PowNumberRule string   `mapstructure:"POW_NUMBER_RULE"`
	NodePorts     []string `mapstructure:"NODE_PORTS"`
}

func NewEnv() (*Env, error) {
	viper.SetConfigFile("../.env")
	if er := viper.ReadInConfig(); er != nil {
		// log.Fatalf("Error reading .env file: %v", err)
		return nil, err.ErrEnvParsing
	}
	viper.AutomaticEnv()

	var env Env
	if er := viper.Unmarshal(&env); er != nil {
		// log.Fatalf("Unable to decode into struct: %v", er)
		return nil, err.ErrEnvParsing
	}

	return &env, nil
}

func (e *Env) GetValueForKey(key string) string {
	return viper.GetString(key)
}
