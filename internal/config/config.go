package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Host         string
	Port         string
	MachineID    int
	Development  bool
	DefaultLimit int
	DBUser       string
	DBPass       string
	DBHost       string
	DBName       string
}

func New(cfgPath string) (*Config, error) {
	var cfg Config
	viper.SetConfigFile(cfgPath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	viper.AutomaticEnv()
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
