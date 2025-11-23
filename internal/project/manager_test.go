package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePath(t *testing.T) {
	manager := NewManager()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "absolute path",
			path:    "/tmp",
			wantErr: false,
		},
		{
			name:    "relative path",
			path:    ".",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := manager.ResolvePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolvePath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
			if err == nil && got == "" {
				t.Errorf("ResolvePath(%q) returned empty path", tt.path)
			}
			if err == nil && !filepath.IsAbs(got) {
				t.Errorf("ResolvePath(%q) returned relative path %q", tt.path, got)
			}
		})
	}
}

func TestResolvePathExpandsHome(t *testing.T) {
	manager := NewManager()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	// Test with ~
	got, err := manager.ResolvePath("~/test")
	if err != nil {
		t.Fatalf("ResolvePath(~/test) error = %v", err)
	}

	expected := filepath.Join(home, "test")
	if got != expected {
		t.Errorf("ResolvePath(~/test) = %q, want %q", got, expected)
	}
}

func TestGlobalConfigPath(t *testing.T) {
	manager := NewManager()

	path, err := manager.GlobalConfigPath()
	if err != nil {
		t.Fatalf("GlobalConfigPath() error = %v", err)
	}

	if path == "" {
		t.Errorf("GlobalConfigPath() returned empty path")
	}

	if !filepath.IsAbs(path) {
		t.Errorf("GlobalConfigPath() returned relative path %q", path)
	}

	if filepath.Base(path) != "config.toml" {
		t.Errorf("GlobalConfigPath() returned path without config.toml filename: %q", path)
	}
}

func TestValidatePath(t *testing.T) {
	manager := NewManager()

	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
	}{
		{
			name: "existing directory",
			setup: func() string {
				return os.TempDir()
			},
			wantErr: false,
		},
		{
			name: "non-existent path",
			setup: func() string {
				return "/nonexistent/path/that/does/not/exist/12345"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			err := manager.ValidatePath(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePath(%q) error = %v, wantErr %v", path, err, tt.wantErr)
			}
		})
	}
}

func TestFormatPathForDisplay(t *testing.T) {
	manager := NewManager()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "path under home directory",
			path: filepath.Join(home, "test", "path"),
			want: "~/test/path",
		},
		{
			name: "absolute path outside home",
			path: "/tmp/test",
			want: "/tmp/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manager.FormatPathForDisplay(tt.path)
			if got != tt.want {
				t.Errorf("FormatPathForDisplay(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
