package config

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// GetString returns specified config in string
func GetString(name string) string {
	return viper.GetString(name)
}

// GetInt returns specified config in int
func GetInt(name string) int {
	return viper.GetInt(name)
}
