package config

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// PeerConfig holds configuration for a peer
type PeerConfig struct {
	Address         string `yaml:"address"`
	TLSRootCertFile string `yaml:"tlsRootCertFile"`
}

// OrdererConfig holds configuration for an orderer
type OrdererConfig struct {
	Address          string `yaml:"address"`
	HostnameOverride string `yaml:"hostnameOverride"`
	TLSCaCert        string `yaml:"tlsCaCert"`
}

// CLIConfig holds configuration for the Fabric CLI
type CLIConfig struct {
	MSP_ID          string `yaml:"mspID"`
	PeerAddress     string `yaml:"peerAddress"`
	TLSRootCertFile string `yaml:"tlsRootCertFile"`
	MSPConfigPath   string `yaml:"mspConfigPath"`
}

// NetworkConfig holds all network-related configuration
type NetworkConfig struct {
	ScriptPath string        `yaml:"scriptPath"`
	CLI        CLIConfig     `yaml:"cli"`
	Orderer    OrdererConfig `yaml:"orderer"`
	Peers      []PeerConfig  `yaml:"peers"`
}

// TimeoutsConfig holds all timeout durations
type TimeoutsConfig struct {
	Network time.Duration `yaml:"network"`
	Deploy  time.Duration `yaml:"deploy"`
	Invoke  time.Duration `yaml:"invoke"`
	Query   time.Duration `yaml:"query"`
	Channel time.Duration `yaml:"channel"`
}

// Config is the top-level configuration structure
type Config struct {
	Network  NetworkConfig  `yaml:"network"`
	Timeouts TimeoutsConfig `yaml:"timeouts"`
}

var (
	cfg  *Config
	once sync.Once
)

// LoadConfig loads configuration from config.yaml, handling defaults and environment variables
func LoadConfig() *Config {
	once.Do(func() {
		log.Println("[INFO] Loading configuration...")
		configPath := os.Getenv("HLF_CONFIG_PATH")
		if configPath == "" {
			configPath = "config.yaml" // Default path
		}

		// Read the config file
		yamlFile, err := os.ReadFile(configPath)
		if err != nil {
			log.Printf("[WARN] Could not read config file '%s': %v. Using default values and environment variables.", configPath, err)
			cfg = loadDefaultAndEnvConfig()
			return
		}

		// Parse the YAML
		var config Config
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			log.Fatalf("[ERROR] Failed to parse config file: %v", err)
		}
		cfg = &config

		// Handle network script path fallback
		if cfg.Network.ScriptPath == "" {
			cfg.Network.ScriptPath = getDefaultScriptPath()
		}

		log.Println("[SUCCESS] Configuration loaded successfully.")
	})
	return cfg
}

// loadDefaultAndEnvConfig provides a fallback if config.yaml is not present
func loadDefaultAndEnvConfig() *Config {
	return &Config{
		Network: NetworkConfig{
			ScriptPath: getDefaultScriptPath(),
			// In a real-world scenario, you might want to load these from env vars as well
			// For now, they are empty and will cause failures if not provided by config.yaml, which is intended.
		},
		Timeouts: TimeoutsConfig{
			Network: 2 * time.Minute,
			Deploy:  5 * time.Minute,
			Invoke:  2 * time.Minute,
			Query:   2 * time.Minute,
			Channel: 2 * time.Minute,
		},
	}
}

// getDefaultScriptPath checks for the network script and fetches it if not found
func getDefaultScriptPath() string {
	path := os.Getenv("HLF_NETWORK_SCRIPT_PATH")
	if path != "" {
		return path
	}

	log.Println("[INFO] HLF_NETWORK_SCRIPT_PATH not provided. Attempting to locate fabric-samples...")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("[ERROR] Unable to determine home directory: %v", err)
	}
	fabricSamplesPath := filepath.Join(homeDir, "fabric-samples")
	networkScriptPath := filepath.Join(fabricSamplesPath, "test-network", "network.sh")

	if _, err := os.Stat(networkScriptPath); os.IsNotExist(err) {
		log.Println("[INFO] fabric-samples not found. Fetching binaries and cloning repository into home directory...")
		cmd := exec.Command("bash", "-c", "curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.4.7 1.5.5")
		cmd.Dir = homeDir // Set the working directory for the download
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("[ERROR] Failed to fetch Fabric binaries and samples: %v", err)
		}
		log.Println("[SUCCESS] Fabric binaries and samples fetched successfully.")
	} else {
		log.Println("[INFO] fabric-samples already exists. Using existing setup.")
	}

	return networkScriptPath
}
