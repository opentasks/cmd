package config

import (
	"bytes"
	_ "embed"
	"fmt"

	"github.com/BurntSushi/toml"
)

//go:embed view.go.tmpl
var ViewTemplate string

// ConfigAsToml returns the configuration formatted as TOML
func ConfigAsToml(cfg *ProjectConfig) (string, error) {
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(cfg); err != nil {
		return "", fmt.Errorf("failed to encode config as TOML: %w", err)
	}
	return buf.String(), nil
}
