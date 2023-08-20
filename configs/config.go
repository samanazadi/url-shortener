package configs

import (
	"github.com/spf13/viper"
)

var Config config

func Init(cfgPath string) error {
	Config = viperConfig{
		fileName: cfgPath,
	}
	return Config.init()
}

type config interface {
	init() error
	GetString(string) string
	GetUint16(string) uint16
	GetBool(string) bool
}

type viperConfig struct {
	fileName string
}

func (v viperConfig) init() error {
	viper.SetConfigFile(v.fileName)
	return viper.ReadInConfig()
}

// GetString returns specified config in string
func (v viperConfig) GetString(name string) string {
	return viper.GetString(name)
}

// GetUint16 returns specified config in int
func (v viperConfig) GetUint16(name string) uint16 {
	return viper.GetUint16(name)
}

func (v viperConfig) GetBool(name string) bool {
	return viper.GetBool(name)
}
