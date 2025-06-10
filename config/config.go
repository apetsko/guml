// Package config provides functionality for managing application configuration.
// It supports loading configuration values from environment variables and command-line flags,
// ensuring flexibility and ease of use in different deployment environments.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/apetsko/guml/utils"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Config string `env:"CONFIG" envDefault:""`

	// Host is the network address with port for the server to listen on.
	Host string `env:"SERVER_ADDRESS" validate:"required"`

	// Cert is the file path to the SSL/TLS certificate used for HTTPS.
	TLSCertPath string `env:"CERT_FILE" validate:"required_if=EnableHTTPS true"`

	// Key is the file path to the SSL/TLS private key used for HTTPS.
	TLSKeyPath string `env:"KEY_FILE" validate:"required_if=EnableHTTPS true"`

	// Https indicates whether the application should use HTTPS for secure communication.
	EnableHTTPS bool `env:"ENABLE_HTTPS"`
}

// New creates a new Config instance, populating it with values from command-line flags and environment variables.
// Returns a pointer to the Config instance or an error if the configuration is invalid.
func New() (*Config, error) {
	var c Config

	// Parse command-line flags
	flag.BoolVar(&c.EnableHTTPS, "s", false, "enable https")
	flag.StringVar(&c.TLSCertPath, "cert", "certs/cert.crt", "certificate filepath")
	flag.StringVar(&c.TLSKeyPath, "key", "certs/cert.key", "private key filepath")
	flag.StringVar(&c.Config, "config", "", "config filepath")
	flag.StringVar(&c.Host, "a", "localhost:8080", "network address with port")

	// Parse the flags
	flag.Parse()

	// Load environment variables into the Config struct
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("failed to load environment: %w", err)
	}

	// Validate the loaded configuration
	if err := utils.ValidateStruct(c); err != nil {
		return nil, err
	}

	// Return the populated and validated Config
	return &c, nil
}

// LoadJSONConfig reads config.json file
func LoadJSONConfig(path string, out interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open config file: %w", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Printf("failed to close config file: %s", err)
		}
	}()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(out); err != nil {
		return fmt.Errorf("decode config: %w", err)
	}

	return nil
}
