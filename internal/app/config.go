// api/config.go

package api

// Config represents the configuration for the API
type Config struct {
	Port        string `json:"bind_port"`
	LoggerLevel string `json:"logger_level"`
}

// NewConfig creates a new instance of the API configuration
func NewConfig() *Config {
	return &Config{
		Port:        ":8080",
		LoggerLevel: "debug",
	}
}