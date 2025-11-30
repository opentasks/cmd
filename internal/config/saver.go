package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// GlobalConfigSaver handles persisting global configuration to disk
type GlobalConfigSaver struct{}

// NewGlobalConfigSaver creates a new GlobalConfigSaver
func NewGlobalConfigSaver() *GlobalConfigSaver {
	return &GlobalConfigSaver{}
}

// Save writes the global config to the specified path
func (s *GlobalConfigSaver) Save(path string, cfg *OpentaskGlobalConfigFile) error {
	data, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// #nosec G306 - File permissions 0600 are intentional for config files (sensitive data)
	// #nosec G304 - Path provided by caller (ProjectManager.GlobalConfigPath), validated
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
