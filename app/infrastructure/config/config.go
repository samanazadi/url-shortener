package config

var configs = make(map[string]any)

func init() {
	configs["server"] = "localhost"
	configs["port"] = "8080"
	configs["dbuser"] = "postgres"
	configs["dbpass"] = "alaki"
	configs["dbhost"] = "localhost"
	configs["dbdb"] = "url-shortener"
}

// GetString returns specified config in string
func GetString(name string) string {
	return configs[name].(string)
}
