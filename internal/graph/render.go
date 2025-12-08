package graph

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/*.dot
var defaultTemplates embed.FS

// RenderDOT generates Graphviz DOT output from a GraphData structure.
// It uses Go templates to render the graph, nodes, and edges.
//
// Parameters:
// - data: The graph data to render
// - templateDir: Directory containing custom templates (empty string uses embedded defaults)
//
// Returns:
// - string: The DOT format output
// - error: Any error encountered during rendering
func RenderDOT(
	data *GraphData,
	templateDir string,
) (string, error) {
	if data == nil {
		return "", fmt.Errorf("data cannot be nil")
	}

	// Load templates
	var tmpl *template.Template
	var err error

	if templateDir != "" {
		// Load custom templates from directory
		tmpl, err = loadCustomTemplates(templateDir)
	} else {
		// Load embedded templates
		tmpl, err = loadEmbeddedTemplates()
	}

	if err != nil {
		return "", fmt.Errorf("failed to load templates: %w", err)
	}

	// Render the graph using the main template
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "graph", data); err != nil {
		return "", fmt.Errorf("failed to render graph: %w", err)
	}

	return buf.String(), nil
}

// loadEmbeddedTemplates loads templates from the embedded filesystem
func loadEmbeddedTemplates() (*template.Template, error) {
	// Parse all templates from embedded filesystem
	tmpl, err := template.ParseFS(defaultTemplates,
		"templates/graph.dot",
		"templates/node.dot",
		"templates/edge.dot",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse embedded templates: %w", err)
	}
	return tmpl, nil
}

// loadCustomTemplates loads templates from a custom directory
func loadCustomTemplates(templateDir string) (*template.Template, error) {
	// Check if directory exists
	info, err := os.Stat(templateDir)
	if err != nil {
		return nil, fmt.Errorf("template directory not found: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("template path is not a directory")
	}

	// Parse all .dot files in the directory
	pattern := filepath.Join(templateDir, "*.dot")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to parse custom templates: %w", err)
	}

	return tmpl, nil
}
