package util

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"gopkg.in/yaml.v3"
)

type Association struct {
	Suffix      string    `yaml:"suffix"`
	Application string    `yaml:"application"`
	BundleID    string    `yaml:"bundle_id"`
	SetAt       time.Time `yaml:"set_at"`
}

type Config struct {
	Version      string                 `yaml:"version"`
	Associations map[string]Association `yaml:"associations"` // key is suffix
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".dutis")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.yaml"), nil
}

func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				Version:      "1.0",
				Associations: make(map[string]Association),
			}, nil
		}
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.Associations == nil {
		config.Associations = make(map[string]Association)
	}

	return &config, nil
}

func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func (c *Config) AddAssociation(suffix, appName, bundleID string) error {
	c.Associations[suffix] = Association{
		Suffix:      suffix,
		Application: appName,
		BundleID:    bundleID,
		SetAt:       time.Now(),
	}
	return c.Save()
}

func (c *Config) RemoveAssociation(suffix string) error {
	delete(c.Associations, suffix)
	return c.Save()
}

func (c *Config) GetAssociation(suffix string) (Association, bool) {
	assoc, ok := c.Associations[suffix]
	return assoc, ok
}

func (c *Config) ListAssociations() []Association {
	var list []Association
	for _, assoc := range c.Associations {
		list = append(list, assoc)
	}
	// Sort by suffix
	sort.Slice(list, func(i, j int) bool {
		return list[i].Suffix < list[j].Suffix
	})
	return list
}

func (c *Config) ApplyAll() error {
	if len(c.Associations) == 0 {
		return fmt.Errorf("no associations configured")
	}

	fmt.Printf("Applying %d file associations...\n\n", len(c.Associations))
	
	successCount := 0
	errorCount := 0

	for _, assoc := range c.ListAssociations() {
		fmt.Printf("  %s → %s (%s)\n", assoc.Suffix, assoc.Application, assoc.BundleID)
		if err := SetDefaultApplication(assoc.BundleID, assoc.Suffix); err != nil {
			fmt.Printf("    ✗ Error: %v\n", err)
			errorCount++
		} else {
			fmt.Printf("    ✓ Applied\n")
			successCount++
		}
	}

	fmt.Printf("\n%d succeeded, %d failed\n", successCount, errorCount)
	
	if errorCount > 0 {
		return fmt.Errorf("%d associations failed to apply", errorCount)
	}
	
	return nil
}
