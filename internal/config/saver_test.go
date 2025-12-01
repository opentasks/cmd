package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestGlobalConfigSaver_Save(t *testing.T) {
	tests := []struct {
		name    string
		config  *OpentaskGlobalConfigFile
		wantErr bool
	}{
		{
			name: "save valid config",
			config: &OpentaskGlobalConfigFile{
				Projects: []GlobalProjectConfig{
					{
						ID:   "test-proj",
						Name: "Test Project",
						Storage: &StorageConfig{
							Backend: "markdown-fs",
							Path:    "/test/path",
						},
						Context: []ProjectContext{
							{Path: "/test/context"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "save empty config",
			config: &OpentaskGlobalConfigFile{
				Projects: []GlobalProjectConfig{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.toml")

			saver := NewGlobalConfigSaver()
			err := saver.Save(configPath, tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file was created
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					t.Error("config file was not created")
					return
				}

				// Verify file contents can be read back
				var readBack OpentaskGlobalConfigFile
				if _, err := toml.DecodeFile(configPath, &readBack); err != nil {
					t.Errorf("failed to read back config: %v", err)
					return
				}

				// Basic validation
				if len(readBack.Projects) != len(tt.config.Projects) {
					t.Errorf("project count mismatch: got %d, want %d", len(readBack.Projects), len(tt.config.Projects))
				}
			}
		})
	}
}

func TestGlobalConfigSaver_SaveInvalidPath(t *testing.T) {
	saver := NewGlobalConfigSaver()
	cfg := &OpentaskGlobalConfigFile{Projects: []GlobalProjectConfig{}}

	// Try to write to non-existent directory without creating it
	invalidPath := "/nonexistent/directory/that/should/not/exist/config.toml"
	err := saver.Save(invalidPath, cfg)

	if err == nil {
		t.Error("expected error when saving to invalid path, got nil")
	}
}
