package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Verbose     bool
	DryRun      bool
	Timeout     int
	EnvVars     map[string]string
	KubeConfig  string
	ProfilePath string
	PresetsPath string
}

func LoadConfig(envFile string) (*Config, error) {
	cfg := &Config{
		EnvVars:     make(map[string]string),
		ProfilePath: "profile",
		PresetsPath: "presets",
	}

	if envFile != "" && fileExists(envFile) {
		if err := godotenv.Load(envFile); err != nil {
			return nil, fmt.Errorf("failed to load env file %s: %w", envFile, err)
		}

		envVars, err := godotenv.Read(envFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read env file %s: %w", envFile, err)
		}

		for key, value := range envVars {
			cfg.EnvVars[key] = value
		}
	}

	cfg.Verbose = viper.GetBool("verbose")
	cfg.DryRun = viper.GetBool("dry-run")
	cfg.Timeout = viper.GetInt("timeout")

	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		cfg.KubeConfig = kubeconfig
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		cfg.KubeConfig = filepath.Join(home, ".kube", "config")
	}

	return cfg, nil
}

func (c *Config) ProcessTemplates(content string) string {
	result := content

	for key, value := range c.EnvVars {
		result = strings.ReplaceAll(result, "${"+key+"}", value)
		result = strings.ReplaceAll(result, "$"+key, value)
	}

	return result
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
