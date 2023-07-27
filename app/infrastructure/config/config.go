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

// GetUint16 returns specified config in int
func GetUint16(name string) uint16 {
	return viper.GetUint16(name)
}
