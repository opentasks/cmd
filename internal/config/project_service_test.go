package config

import (
	"testing"
)

func TestProjectService_AttachContext(t *testing.T) {
	tests := []struct {
		name          string
		initialConfig *OpentaskGlobalConfigFile
		projectID     string
		path          string
		wantErr       bool
		wantContexts  int
	}{
		{
			name: "attach to existing project",
			initialConfig: &OpentaskGlobalConfigFile{
				Projects: []GlobalProjectConfig{
					{ID: "test-project", Name: "Test", Context: []ProjectContext{}},
				},
			},
			projectID:    "test-project",
			path:         "/path/to/context",
			wantErr:      false,
			wantContexts: 1,
		},
		{
			name: "attach to non-existent project",
			initialConfig: &OpentaskGlobalConfigFile{
				Projects: []GlobalProjectConfig{},
			},
			projectID:    "missing-project",
			path:         "/path/to/context",
			wantErr:      true,
			wantContexts: 0,
		},
		{
			name: "attach duplicate path (returns error)",
			initialConfig: &OpentaskGlobalConfigFile{
				Projects: []GlobalProjectConfig{
					{
						ID:   "test-project",
						Name: "Test",
						Context: []ProjectContext{
							{Path: "/existing/path"},
						},
					},
				},
			},
			projectID:    "test-project",
			path:         "/existing/path",
			wantErr:      true,
			wantContexts: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProjectService(tt.initialConfig)
			err := svc.AttachContext(tt.projectID, tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("AttachContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				proj := svc.GetProject(tt.projectID)
				if proj == nil {
					t.Fatal("project not found after attach")
				}
				if len(proj.Context) != tt.wantContexts {
					t.Errorf("got %d contexts, want %d", len(proj.Context), tt.wantContexts)
				}
			}
		})
	}
}

func TestProjectService_DetachContext(t *testing.T) {
	tests := []struct {
		name          string
		initialConfig *OpentaskGlobalConfigFile
		projectID     string
		path          string
		wantErr       bool
		wantContexts  int
	}{
		{
			name: "detach existing context",
			initialConfig: &OpentaskGlobalConfigFile{
				Projects: []GlobalProjectConfig{
					{
						ID:   "test-project",
						Name: "Test",
						Context: []ProjectContext{
							{Path: "/path/to/remove"},
							{Path: "/path/to/keep"},
						},
					},
				},
			},
			projectID:    "test-project",
			path:         "/path/to/remove",
			wantErr:      false,
			wantContexts: 1,
		},
		{
			name: "detach from non-existent project",
			initialConfig: &OpentaskGlobalConfigFile{
				Projects: []GlobalProjectConfig{},
			},
			projectID:    "missing-project",
			path:         "/path/to/context",
			wantErr:      true,
			wantContexts: 0,
		},
		{
			name: "detach non-existent path (returns error)",
			initialConfig: &OpentaskGlobalConfigFile{
				Projects: []GlobalProjectConfig{
					{
						ID:      "test-project",
						Name:    "Test",
						Context: []ProjectContext{{Path: "/existing"}},
					},
				},
			},
			projectID:    "test-project",
			path:         "/not-there",
			wantErr:      true,
			wantContexts: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProjectService(tt.initialConfig)
			err := svc.DetachContext(tt.projectID, tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("DetachContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				proj := svc.GetProject(tt.projectID)
				if proj == nil {
					t.Fatal("project not found after detach")
				}
				if len(proj.Context) != tt.wantContexts {
					t.Errorf("got %d contexts, want %d", len(proj.Context), tt.wantContexts)
				}
			}
		})
	}
}

func TestProjectService_RemoveProject(t *testing.T) {
	tests := []struct {
		name              string
		initialConfig     *OpentaskGlobalConfigFile
		projectID         string
		wantErr           bool
		wantProjectCount  int
		wantActiveCleared bool
	}{
		{
			name: "remove existing project",
			initialConfig: &OpentaskGlobalConfigFile{
				Projects: []GlobalProjectConfig{
					{ID: "proj1", Name: "Project 1"},
					{ID: "proj2", Name: "Project 2"},
				},
				ActiveProject: "proj1",
			},
			projectID:         "proj2",
			wantErr:           false,
			wantProjectCount:  1,
			wantActiveCleared: false,
		},
		{
			name: "remove active project (should clear active_project)",
			initialConfig: &OpentaskGlobalConfigFile{
				Projects: []GlobalProjectConfig{
					{ID: "proj1", Name: "Project 1"},
					{ID: "proj2", Name: "Project 2"},
				},
				ActiveProject: "proj1",
			},
			projectID:         "proj1",
			wantErr:           false,
			wantProjectCount:  1,
			wantActiveCleared: true,
		},
		{
			name: "remove non-existent project",
			initialConfig: &OpentaskGlobalConfigFile{
				Projects: []GlobalProjectConfig{
					{ID: "proj1", Name: "Project 1"},
				},
			},
			projectID:        "missing",
			wantErr:          true,
			wantProjectCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProjectService(tt.initialConfig)
			err := svc.RemoveProject(tt.projectID)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(tt.initialConfig.Projects) != tt.wantProjectCount {
					t.Errorf("got %d projects, want %d", len(tt.initialConfig.Projects), tt.wantProjectCount)
				}

				if tt.wantActiveCleared && tt.initialConfig.ActiveProject != "" {
					t.Errorf("active_project not cleared, got %s", tt.initialConfig.ActiveProject)
				}

				// Verify project is actually removed
				if svc.GetProject(tt.projectID) != nil {
					t.Error("project still exists after removal")
				}
			}
		})
	}
}

func TestProjectService_GetProject(t *testing.T) {
	cfg := &OpentaskGlobalConfigFile{
		Projects: []GlobalProjectConfig{
			{ID: "proj1", Name: "Project 1"},
			{ID: "proj2", Name: "Project 2"},
		},
	}

	svc := NewProjectService(cfg)

	t.Run("get existing project", func(t *testing.T) {
		proj := svc.GetProject("proj1")
		if proj == nil {
			t.Fatal("expected project, got nil")
		}
		if proj.ID != "proj1" {
			t.Errorf("got ID %s, want proj1", proj.ID)
		}
	})

	t.Run("get non-existent project", func(t *testing.T) {
		proj := svc.GetProject("missing")
		if proj != nil {
			t.Error("expected nil, got project")
		}
	})
}
