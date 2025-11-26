// Package brand provides branding assets for the opentask CLI.
package brand

import (
	"embed"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mbndr/figlet4go"
)

//go:embed fonts/*.flf
var fontsFS embed.FS

// Version can be set at build time
var Version = "dev"

// AppName is the name rendered in the banner
const AppName = "opentask"

// availableFonts lists fonts that work well for the banner
var availableFonts = []string{
	"slant",
	"small",
	"smslant",
	"doom",
	"big",
	"shadow",
	"mini",
	"script",
}

// Theme represents a color theme for the banner
type Theme struct {
	Name   string
	Accent string // first letter color
	Base   string // remaining letters color
}

// themes is the collection of available banner themes
var themes = []Theme{
	{Name: "electric", Accent: "#06b6d4", Base: "#94a3b8"},   // cyan accent, slate base
	{Name: "ember", Accent: "#f97316", Base: "#fafaf9"},      // orange accent, warm white
	{Name: "grape", Accent: "#a855f7", Base: "#e2e8f0"},      // purple accent, light slate
	{Name: "mint", Accent: "#10b981", Base: "#f0fdf4"},       // emerald accent, green-white
	{Name: "rose", Accent: "#f43f5e", Base: "#fafafa"},       // rose accent, white
	{Name: "gold", Accent: "#eab308", Base: "#fefce8"},       // yellow accent, cream
	{Name: "midnight", Accent: "#3b82f6", Base: "#94a3b8"},   // blue accent, slate
	{Name: "neon", Accent: "#22d3ee", Base: "#334155"},       // bright cyan, dark slate
	{Name: "blood", Accent: "#dc2626", Base: "#d6d3d1"},      // red accent, stone
	{Name: "hacker", Accent: "#22c55e", Base: "#166534"},     // bright green, dark green
	{Name: "sunset", Accent: "#fb923c", Base: "#fef3c7"},     // orange, amber-white
	{Name: "ocean", Accent: "#0ea5e9", Base: "#e0f2fe"},      // sky blue, light blue
	{Name: "synthwave", Accent: "#e879f9", Base: "#c084fc"},  // fuchsia, purple
	{Name: "monochrome", Accent: "#f8fafc", Base: "#64748b"}, // white accent, slate base
}

// renderer is the lipgloss renderer for stdout (for tagline)
var renderer = lipgloss.NewRenderer(os.Stdout)

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

		// Extract font name without extension
		fontName := strings.TrimSuffix(entry.Name(), ".flf")
		_ = figletRenderer.LoadBindataFont(data, fontName)
	}
}

// randomTheme returns a random theme
func randomTheme() Theme {
	return themes[rand.Intn(len(themes))]
}

// randomFont returns a random font name
func randomFont() string {
	return availableFonts[rand.Intn(len(availableFonts))]
}

// hexToColor converts a hex color string to figlet4go.Color
func hexToColor(hex string) figlet4go.Color {
	// Strip # prefix if present
	hex = strings.TrimPrefix(hex, "#")
	color, err := figlet4go.NewTrueColorFromHexString(hex)
	if err != nil {
		// Fallback to white on error
		color, _ = figlet4go.NewTrueColorFromHexString("ffffff")
	}
	return color
}

// buildColorSlice creates a color slice for figlet4go with accent on first char
func buildColorSlice(theme Theme, textLen int) []figlet4go.Color {
	colors := make([]figlet4go.Color, textLen)
	accent := hexToColor(theme.Accent)
	base := hexToColor(theme.Base)

	for i := 0; i < textLen; i++ {
		if i == 0 {
			colors[i] = accent
		} else {
			colors[i] = base
		}
	}
	return colors
}

// renderBanner generates the ASCII art banner with theme colors and random font
func renderBanner(theme Theme, font string) string {
	options := figlet4go.NewRenderOptions()
	options.FontName = font
	options.FontColor = buildColorSlice(theme, len(AppName))

	result, err := figletRenderer.RenderOpts(AppName, options)
	if err != nil {
		// Fallback to plain text on error
		return AppName + "\n"
	}
	return result
}

// PrintBanner writes the banner to the given writer
func PrintBanner(w io.Writer) {
	theme := randomTheme()
	font := randomFont()
	fmt.Fprint(w, renderBanner(theme, font))
}

// PrintBannerWithVersion writes the banner with version info to the given writer
func PrintBannerWithVersion(w io.Writer) {
	theme := randomTheme()
	font := randomFont()
	fmt.Fprint(w, renderBanner(theme, font))

	// Style the tagline with the accent color
	taglineStyle := renderer.NewStyle().
		Foreground(lipgloss.Color(theme.Accent)).
		Faint(true)
	fmt.Fprintf(w, "  %s\n\n", taglineStyle.Render(fmt.Sprintf("Task management with markdown files (v%s)", Version)))
}
