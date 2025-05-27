package config

import (
	"os"
)

// Config holds configuration values for the controller
type Config struct {
	NetworkScriptPath string
}

// LoadConfig loads configuration from environment variables or uses defaults
func LoadConfig() *Config {
	path := os.Getenv("HLF_NETWORK_SCRIPT_PATH")
	if path == "" {
		path = "/Users/padamarajkore/fabric-samples/test-network/network.sh"
	}
	return &Config{
		NetworkScriptPath: path,
	}
}
