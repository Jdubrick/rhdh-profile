package kustomize

import (
	"fmt"
	"os"
	"strings"

	"github.com/Jdubrick/rhdh-profile/pkg/config"
)

type Kustomizer struct {
	config *config.Config
}

func New(cfg *config.Config) (*Kustomizer, error) {
	return &Kustomizer{
		config: cfg,
	}, nil
}

func (k *Kustomizer) UpdateSecretsFile(secretsPath string, envVars map[string]string) error {
	if k.config.Verbose {
		fmt.Printf("Updating secrets file: %s with %d environment variables\n", secretsPath, len(envVars))
	}

	content, err := os.ReadFile(secretsPath)
	if err != nil {
		return fmt.Errorf("failed to read secrets file %s: %w", secretsPath, err)
	}

	secretsContent := string(content)

	for envKey, envValue := range envVars {
		// Look for the key in the stringData section (handles both quoted and template formats)
		keyPattern := fmt.Sprintf(`%s:`, envKey)
		if strings.Contains(secretsContent, keyPattern) {
			lines := strings.Split(secretsContent, "\n")
			for i, line := range lines {
				// Check if this line contains our key and it's the exact key (not a substring)
				if strings.Contains(line, keyPattern) {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) == 2 && strings.TrimSpace(parts[0]) == envKey {
						// Extract indentation
						trimmed := strings.TrimLeft(line, " \t")
						indent := line[:len(line)-len(trimmed)]

						// Handle multi-line values (like private keys) using literal block scalar
						if strings.Contains(envValue, "\n") {
							// Remove any existing multi-line content for this key first
							endLine := i
							for j := i + 1; j < len(lines); j++ {
								// Check if this line is still part of the multi-line value
								if strings.HasPrefix(lines[j], indent+"  ") || strings.TrimSpace(lines[j]) == "" {
									endLine = j
								} else {
									break
								}
							}

							// Remove the old multi-line content
							if endLine > i {
								lines = append(lines[:i+1], lines[endLine+1:]...)
							}

							// Use YAML literal block scalar for multi-line values
							lines[i] = fmt.Sprintf(`%s%s: |`, indent, envKey)
							// Add each line of the multi-line value with proper indentation
							valueLines := strings.Split(envValue, "\n")
							for j := len(valueLines) - 1; j >= 0; j-- {
								lines = append(lines[:i+1], append([]string{indent + "  " + valueLines[j]}, lines[i+1:]...)...)
							}
						} else {
							// Single line value - use quotes for safety
							lines[i] = fmt.Sprintf(`%s%s: "%s"`, indent, envKey, envValue)
						}

						if k.config.Verbose {
							fmt.Printf("Updated %s in secrets file\n", envKey)
						}
						break
					}
				}
			}
			secretsContent = strings.Join(lines, "\n")
		} else if k.config.Verbose {
			fmt.Printf("Key %s not found in secrets file, skipping\n", envKey)
		}
	}

	if err := os.WriteFile(secretsPath, []byte(secretsContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated secrets file: %w", err)
	}

	if k.config.Verbose {
		fmt.Printf("Successfully updated secrets file: %s\n", secretsPath)
	}

	return nil
}
