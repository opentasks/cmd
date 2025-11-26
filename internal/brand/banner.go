// Package brand provides branding assets for the opentask CLI.
package brand

import (
	"embed"
	"fmt"
	"io"
	"math/rand"
	"path/filepath"
	"strings"

	"github.com/mbndr/figlet4go"
)

//go:embed fonts/*.flf
var fontsFS embed.FS

// Version can be set at build time
var Version = "dev"

// figletRenderer is the shared figlet renderer with embedded fonts loaded
var figletRenderer *figlet4go.AsciiRender

func init() {
	figletRenderer = figlet4go.NewAsciiRender()
	loadEmbeddedFonts()
}

// loadEmbeddedFonts loads all .flf fonts from the embedded filesystem
func loadEmbeddedFonts() {
	entries, err := fontsFS.ReadDir("fonts")
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".flf") {
			continue
		}

		data, err := fontsFS.ReadFile(filepath.Join("fonts", entry.Name()))
		if err != nil {
			continue
		}

		fontName := strings.TrimSuffix(entry.Name(), ".flf")
		_ = figletRenderer.LoadBindataFont(data, fontName)
	}
}

// Banner renders ASCII art text with configurable themes and fonts
type Banner struct {
	text    string
	theme   Theme
	font    string
	tagline string
	version string
}

// Option is a functional option for configuring a Banner
type Option func(*Banner)

// WithText sets the text to render
func WithText(text string) Option {
	return func(b *Banner) {
		b.text = text
	}
}

// WithTheme sets the color theme
func WithTheme(theme Theme) Option {
	return func(b *Banner) {
		b.theme = theme
	}
}

// WithFont sets the FIGlet font
func WithFont(font string) Option {
	return func(b *Banner) {
		b.font = font
	}
}

// WithTagline sets the tagline displayed below the banner
func WithTagline(tagline string) Option {
	return func(b *Banner) {
		b.tagline = tagline
	}
}

// WithVersion sets the version displayed in the tagline
func WithVersion(version string) Option {
	return func(b *Banner) {
		b.version = version
	}
}

// WithRandomTheme selects a random theme
func WithRandomTheme() Option {
	return func(b *Banner) {
		b.theme = Themes[rand.Intn(len(Themes))]
	}
}

// WithRandomFont selects a random font
func WithRandomFont() Option {
	return func(b *Banner) {
		b.font = Fonts[rand.Intn(len(Fonts))]
	}
}

// NewBanner creates a new Banner with the given options
func NewBanner(opts ...Option) *Banner {
	b := &Banner{
		text:    "banner",
		theme:   Themes[0], // Default to first curated theme
		font:    "slant",
		version: Version,
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

// renderFiglet generates plain ASCII art without colors
func (b *Banner) renderFiglet() ([]string, error) {
	options := figlet4go.NewRenderOptions()
	options.FontName = b.font

	result, err := figletRenderer.RenderOpts(b.text, options)
	if err != nil {
		return nil, err
	}

	// Split into lines, removing trailing empty lines
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	return lines, nil
}

// applyTheme applies lipgloss styles to figlet output
// First character column gets accent style, rest gets base style
func (b *Banner) applyTheme(lines []string) string {
	if len(lines) == 0 {
		return ""
	}

	// Find character column boundaries by analyzing the figlet output
	// Each character in the input text corresponds to a "column" in the output
	charWidths := b.findCharWidths(lines)

	var result strings.Builder
	for _, line := range lines {
		styledLine := b.styleLineByColumns(line, charWidths)
		result.WriteString(styledLine)
		result.WriteString("\n")
	}

	return result.String()
}

// findCharWidths analyzes figlet output to find width of each character column
func (b *Banner) findCharWidths(lines []string) []int {
	if len(lines) == 0 || len(b.text) == 0 {
		return nil
	}

	// Render each character individually to find its width
	widths := make([]int, len(b.text))
	for i, char := range b.text {
		options := figlet4go.NewRenderOptions()
		options.FontName = b.font
		charResult, err := figletRenderer.RenderOpts(string(char), options)
		if err != nil {
			widths[i] = 1
			continue
		}
		charLines := strings.Split(charResult, "\n")
		if len(charLines) > 0 {
			widths[i] = len(charLines[0])
		}
	}

	return widths
}

// styleLineByColumns applies accent to first char column, base to rest
func (b *Banner) styleLineByColumns(line string, charWidths []int) string {
	if len(charWidths) == 0 || len(line) == 0 {
		return b.theme.Base.Render(line)
	}

	firstCharWidth := charWidths[0]
	if firstCharWidth >= len(line) {
		return b.theme.Accent.Render(line)
	}

	accent := line[:firstCharWidth]
	base := line[firstCharWidth:]

	return b.theme.Accent.Render(accent) + b.theme.Base.Render(base)
}

// Render generates the ASCII art banner string with theme applied
func (b *Banner) Render() string {
	lines, err := b.renderFiglet()
	if err != nil {
		return b.theme.Accent.Render(b.text) + "\n"
	}

	return b.applyTheme(lines)
}

// RenderWithTagline generates the banner with a styled tagline
func (b *Banner) RenderWithTagline() string {
	result := b.Render()

	if b.tagline != "" {
		taglineText := b.tagline
		if b.version != "" {
			taglineText = fmt.Sprintf("%s (v%s)", b.tagline, b.version)
		}
		result += fmt.Sprintf("  %s\n\n", b.theme.Tagline.Render(taglineText))
	}

	return result
}

// Print writes the banner to the given writer
func (b *Banner) Print(w io.Writer) {
	fmt.Fprint(w, b.Render())
}

// PrintWithTagline writes the banner with tagline to the given writer
func (b *Banner) PrintWithTagline(w io.Writer) {
	fmt.Fprint(w, b.RenderWithTagline())
}

// --- Convenience functions for backward compatibility ---

// PrintBanner writes a random banner to the given writer
func PrintBanner(w io.Writer) {
	b := NewBanner(
		WithText("opentask"),
		WithRandomTheme(),
		WithRandomFont(),
	)
	b.Print(w)
}

// PrintBannerWithVersion writes a random banner with version to the given writer
func PrintBannerWithVersion(w io.Writer) {
	b := NewBanner(
		WithText("opentask"),
		WithRandomTheme(),
		WithRandomFont(),
		WithTagline("Task management with markdown files"),
		WithVersion(Version),
	)
	b.PrintWithTagline(w)
}
