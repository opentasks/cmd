package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/glamour"
	"github.com/zenobi-us/opentask/internal/display"
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

type RendererOptions struct {
	Resolved *OpentaskResolvedConfig
	Verbose  bool
	Cwd      string
}

func WithVerbose(verbose bool) func(*RendererOptions) {
	return func(ro *RendererOptions) {
		ro.Verbose = verbose
	}
}

func WithCwd(cwd string) func(*RendererOptions) {
	return func(ro *RendererOptions) {
		ro.Cwd = cwd
	}
}

func WithResolvedConfig(resolved *OpentaskResolvedConfig) func(*RendererOptions) {
	return func(ro *RendererOptions) {
		ro.Resolved = resolved
	}
}

func Renderer(opts ...func(*RendererOptions)) error {
	options := &RendererOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Build file tree
	fileTree := display.ConfigFileTree(options.Resolved.DiscoveredFiles)

	// Convert resolved config to legacy ProjectConfig format for display
	displayConfig := &ProjectConfig{
		Project:   *options.Resolved.Project,
		Workflow:  *options.Resolved.Workflow,
		Templates: *options.Resolved.Templates,
		Storage:   *options.Resolved.Storage,
	}

	// Convert config to TOML string
	tomlStr, err := ConfigAsToml(displayConfig)
	if err != nil {
		return fmt.Errorf("failed to format config as TOML: %w", err)
	}

	// Prepare template data
	templateData := map[string]interface{}{
		"FoundFiles":                 options.Resolved.DiscoveredFiles,
		"MergingOrder":               options.Resolved.DiscoveredFiles,
		"ResolvedConfigAsTomlString": tomlStr,
		"StopReason":                 "reached filesystem root",
		"StopDir":                    options.Cwd,
		"FileTree":                   fileTree,
	}

	// Execute template
	tmpl, err := template.New("configView").
		Funcs(template.FuncMap{
			"quote": func(s string) string {
				return `"` + s + `"`
			},
			"add": func(a, b int) int {
				return a + b
			},
		}).
		Parse(ViewTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Render markdown to a buffer first
	var mdBuf bytes.Buffer
	if err := tmpl.Execute(&mdBuf, templateData); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Render markdown with syntax highlighting using glamour
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)
	if err != nil {
		return fmt.Errorf("failed to create markdown renderer: %w", err)
	}

	output, err := renderer.Render(mdBuf.String())
	if err != nil {
		return fmt.Errorf("failed to render markdown: %w", err)
	}

	fmt.Print(output)

	// Verbose output
	if options.Verbose {
		fmt.Println("\n=== Verbose Mode: Merging Details ===")

		// Show discovered config files
		if len(options.Resolved.DiscoveredFiles) > 0 {
			for i, file := range options.Resolved.DiscoveredFiles {
				fmt.Printf("\n[Step %d] Applying: %s\n", i+1, file)
				// For resolved config, just show what was merged
				if options.Resolved.Project.Name != "" {
					fmt.Printf("  - project.name: %q\n", options.Resolved.Project.Name)
				}
				if len(options.Resolved.Workflow.Statuses) > 0 {
					fmt.Printf("  - workflow.statuses: %v\n", options.Resolved.Workflow.Statuses)
				}
				if options.Resolved.Storage.Path != "" {
					fmt.Printf("  - storage.path: %q\n", options.Resolved.Storage.Path)
				}
			}
		}

		// Show defaults as final virtual layer
		stepNum := len(options.Resolved.DiscoveredFiles) + 1
		fmt.Printf("\n[Step %d] (Virtual) Default configuration\n", stepNum)
		defaults := ProjectConfig{
			Workflow:  DefaultWorkflow(),
			Templates: DefaultTemplates(),
			Storage:   DefaultStorage(),
		}
		fmt.Printf("  - workflow.statuses: %v\n", defaults.Workflow.Statuses)
		fmt.Printf("  - workflow.initial: %q\n", defaults.Workflow.Initial)
		fmt.Printf("  - storage.backend: %q\n", defaults.Storage.Backend)
	}

	return nil
}
