package project

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"
)

//go:embed init.md.tmpl
var initTemplate string

// InitialTemplate returns the initial AGENTS guidelines template
var InitialTemplate = initTemplate

// RenderAgentsGuide renders the AGENTS guidelines template with the given context data
func RenderAgentsGuide(data map[string]string) (string, error) {
	tmpl, err := template.New("agents-guide").Parse(InitialTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return result.String(), nil
}
