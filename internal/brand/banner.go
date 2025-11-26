// Package brand provides branding assets for the opentask CLI.
package brand

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Version can be set at build time
var Version = "dev"

// ThemeType distinguishes between gradient and duotone themes
type ThemeType int

const (
	ThemeGradient ThemeType = iota
	ThemeDuotone
)

// Theme represents a color theme for the banner
type Theme struct {
	Name   string
	Type   ThemeType
	Colors []string // gradient: top to bottom, duotone: [accent, base]
}

// themes is the collection of available banner themes
var themes = []Theme{
	// Gradient themes
	{
		Name:   "ocean",
		Type:   ThemeGradient,
		Colors: []string{"#0ea5e9", "#06b6d4", "#14b8a6", "#10b981", "#22c55e", "#84cc16"},
	},
	{
		Name:   "sunset",
		Type:   ThemeGradient,
		Colors: []string{"#f97316", "#fb923c", "#f59e0b", "#eab308", "#facc15", "#fde047"},
	},
	{
		Name:   "aurora",
		Type:   ThemeGradient,
		Colors: []string{"#a855f7", "#8b5cf6", "#6366f1", "#3b82f6", "#0ea5e9", "#06b6d4"},
	},
	{
		Name:   "forest",
		Type:   ThemeGradient,
		Colors: []string{"#166534", "#15803d", "#16a34a", "#22c55e", "#4ade80", "#86efac"},
	},
	{
		Name:   "fire",
		Type:   ThemeGradient,
		Colors: []string{"#dc2626", "#ef4444", "#f97316", "#fb923c", "#fbbf24", "#fde047"},
	},
	{
		Name:   "synthwave",
		Type:   ThemeGradient,
		Colors: []string{"#f472b6", "#e879f9", "#c084fc", "#a78bfa", "#818cf8", "#60a5fa"},
	},
	{
		Name:   "monochrome",
		Type:   ThemeGradient,
		Colors: []string{"#f8fafc", "#e2e8f0", "#cbd5e1", "#94a3b8", "#64748b", "#475569"},
	},
	{
		Name:   "candy",
		Type:   ThemeGradient,
		Colors: []string{"#fb7185", "#f472b6", "#e879f9", "#c084fc", "#a78bfa", "#818cf8"},
	},
	{
		Name:   "matrix",
		Type:   ThemeGradient,
		Colors: []string{"#4ade80", "#22c55e", "#16a34a", "#15803d", "#166534", "#14532d"},
	},
	{
		Name:   "ice",
		Type:   ThemeGradient,
		Colors: []string{"#e0f2fe", "#bae6fd", "#7dd3fc", "#38bdf8", "#0ea5e9", "#0284c7"},
	},
	// Duotone themes [accent, base]
	{
		Name:   "duotone-ember",
		Type:   ThemeDuotone,
		Colors: []string{"#f97316", "#fafaf9"}, // orange accent, warm white base
	},
	{
		Name:   "duotone-electric",
		Type:   ThemeDuotone,
		Colors: []string{"#06b6d4", "#f1f5f9"}, // cyan accent, cool white base
	},
	{
		Name:   "duotone-grape",
		Type:   ThemeDuotone,
		Colors: []string{"#a855f7", "#e2e8f0"}, // purple accent, slate base
	},
	{
		Name:   "duotone-mint",
		Type:   ThemeDuotone,
		Colors: []string{"#10b981", "#f0fdf4"}, // emerald accent, green-tinted white
	},
	{
		Name:   "duotone-rose",
		Type:   ThemeDuotone,
		Colors: []string{"#f43f5e", "#fafafa"}, // rose accent, neutral white
	},
	{
		Name:   "duotone-gold",
		Type:   ThemeDuotone,
		Colors: []string{"#eab308", "#fefce8"}, // yellow accent, warm cream base
	},
	{
		Name:   "duotone-midnight",
		Type:   ThemeDuotone,
		Colors: []string{"#3b82f6", "#94a3b8"}, // blue accent, slate base
	},
	{
		Name:   "duotone-neon",
		Type:   ThemeDuotone,
		Colors: []string{"#22d3ee", "#0f172a"}, // bright cyan accent, dark slate base
	},
	{
		Name:   "duotone-blood",
		Type:   ThemeDuotone,
		Colors: []string{"#dc2626", "#1c1917"}, // red accent, dark stone base
	},
	{
		Name:   "duotone-hacker",
		Type:   ThemeDuotone,
		Colors: []string{"#22c55e", "#052e16"}, // green accent, dark green base
	},
}

// renderer is the lipgloss renderer for stdout
var renderer = lipgloss.NewRenderer(os.Stdout)

// randomTheme returns a random theme
func randomTheme() Theme {
	return themes[rand.Intn(len(themes))]
}

// bannerLines are the raw ASCII art lines
var bannerLines = []string{
	`                         __            __  `,
	`  ____  ____  ___  ____  / /_____ ____ / /__`,
	` / __ \/ __ \/ _ \/ __ \/ __/ __ '/ __/ //_/`,
	`/ /_/ / /_/ /  __/ / / / /_/ /_/ (__  )  <  `,
	`\____/ .___/\___/_/ /_/\__/\__,_/____/_/|_| `,
	`    /_/                                     `,
}

// firstCharIndex returns the index of the first non-space character
func firstCharIndex(s string) int {
	for i, r := range s {
		if r != ' ' {
			return i
		}
	}
	return -1
}

// renderGradient renders the banner with gradient colors
func renderGradient(theme Theme) string {
	var result string
	for i, line := range bannerLines {
		colorIdx := i
		if colorIdx >= len(theme.Colors) {
			colorIdx = len(theme.Colors) - 1
		}
		style := renderer.NewStyle().Foreground(lipgloss.Color(theme.Colors[colorIdx]))
		result += style.Render(line) + "\n"
	}
	return result
}

// renderDuotone renders the banner with first-char accent and rest in base color
func renderDuotone(theme Theme) string {
	accent := renderer.NewStyle().Foreground(lipgloss.Color(theme.Colors[0])).Bold(true)
	base := renderer.NewStyle().Foreground(lipgloss.Color(theme.Colors[1]))

	var result string
	for _, line := range bannerLines {
		idx := firstCharIndex(line)
		if idx == -1 {
			// All spaces
			result += line + "\n"
			continue
		}

		// Build: leading spaces + accent char + rest in base
		var sb strings.Builder
		sb.WriteString(line[:idx])                       // leading spaces
		sb.WriteString(accent.Render(string(line[idx]))) // first char accented
		sb.WriteString(base.Render(line[idx+1:]))        // rest in base
		result += sb.String() + "\n"
	}
	return result
}

// renderBanner applies a theme to the banner text
func renderBanner(theme Theme) string {
	switch theme.Type {
	case ThemeDuotone:
		return renderDuotone(theme)
	default:
		return renderGradient(theme)
	}
}

// PrintBanner writes the banner to the given writer
func PrintBanner(w io.Writer) {
	theme := randomTheme()
	fmt.Fprint(w, renderBanner(theme))
}

// PrintBannerWithVersion writes the banner with version info to the given writer
func PrintBannerWithVersion(w io.Writer) {
	theme := randomTheme()
	fmt.Fprint(w, renderBanner(theme))

	// Style the tagline
	var taglineColor string
	if theme.Type == ThemeDuotone {
		taglineColor = theme.Colors[0] // use accent for tagline
	} else {
		taglineColor = theme.Colors[len(theme.Colors)-1]
	}
	taglineStyle := renderer.NewStyle().
		Foreground(lipgloss.Color(taglineColor)).
		Faint(true)
	fmt.Fprintf(w, "  %s\n\n", taglineStyle.Render(fmt.Sprintf("Task management with markdown files (v%s)", Version)))
}
