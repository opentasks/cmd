package brand

import (
	"math/rand"

	"github.com/charmbracelet/lipgloss"
	goghthemes "github.com/willyv3/gogh-themes/lipgloss"
)

// Theme represents a color theme for the banner
type Theme struct {
	Name    string
	Accent  lipgloss.Style // first letter style
	Base    lipgloss.Style // remaining letters style
	Tagline lipgloss.Style // tagline style
}

// ThemeFromGogh creates a Theme from a gogh-themes color scheme
func ThemeFromGogh(name string) Theme {
	gogh, ok := goghthemes.Get(name)
	if !ok {
		// Fallback to a default
		gogh, _ = goghthemes.Get("Dracula")
	}

	return Theme{
		Name:    name,
		Accent:  lipgloss.NewStyle().Foreground(gogh.BrightCyan).Bold(true),
		Base:    lipgloss.NewStyle().Foreground(gogh.Foreground),
		Tagline: lipgloss.NewStyle().Foreground(gogh.BrightCyan).Faint(true),
	}
}

// ThemeFromGoghWithAccent creates a Theme using a specific accent color from the palette
func ThemeFromGoghWithAccent(name string, accentColor func(goghthemes.Theme) lipgloss.Color) Theme {
	gogh, ok := goghthemes.Get(name)
	if !ok {
		gogh, _ = goghthemes.Get("Dracula")
	}

	accent := accentColor(gogh)
	return Theme{
		Name:    name,
		Accent:  lipgloss.NewStyle().Foreground(accent).Bold(true),
		Base:    lipgloss.NewStyle().Foreground(gogh.Foreground),
		Tagline: lipgloss.NewStyle().Foreground(accent).Faint(true),
	}
}

// Curated theme names that look good for banners
var curatedThemeNames = []string{
	"Dracula",
	"Nord",
	"Gruvbox Dark",
	"Catppuccin Mocha",
	"Catppuccin Frappe",
	"Tokyo Night",
	"Tokyo Night Storm",
	"Solarized Dark",
	"One Dark",
	"Monokai Pro",
	"Ayu Dark",
	"Ayu Mirage",
	"Nightfox",
	"Kanagawa",
	"Rose Pine",
	"Rose Pine Moon",
	"Everforest Dark Hard",
	"Cyberpunk",
	"Synthwave 84",
	"Shades Of Purple",
}

// Themes is the collection of curated banner themes
var Themes []Theme

func init() {
	// Build themes from curated list
	Themes = make([]Theme, 0, len(curatedThemeNames))
	for _, name := range curatedThemeNames {
		if _, ok := goghthemes.Get(name); ok {
			Themes = append(Themes, ThemeFromGogh(name))
		}
	}

	// Fallback if none found
	if len(Themes) == 0 {
		Themes = append(Themes, ThemeFromGogh("Dracula"))
	}
}

// RandomTheme returns a random theme from the curated list
func RandomTheme() Theme {
	return Themes[rand.Intn(len(Themes))]
}

// AllGoghThemes returns all 361 gogh themes as banner Themes
func AllGoghThemes() []Theme {
	names := goghthemes.Names()
	themes := make([]Theme, len(names))
	for i, name := range names {
		themes[i] = ThemeFromGogh(name)
	}
	return themes
}

// GetTheme returns a specific theme by name
func GetTheme(name string) (Theme, bool) {
	_, ok := goghthemes.Get(name)
	if !ok {
		return Theme{}, false
	}
	return ThemeFromGogh(name), true
}

// Fonts lists available FIGlet fonts
var Fonts = []string{
	"slant",
	"small",
	"smslant",
	"doom",
	"big",
	"shadow",
	"mini",
	"script",
}
