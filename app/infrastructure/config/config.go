package config

var configs map[string]any

func init() {
	configs["server"] = "localhost"
	configs["port"] = "8080"
}

// GetString returns specified config in string
func GetString(name string) string {
	return configs[name].(string)
}
