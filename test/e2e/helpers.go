package e2e

import (
	"context"
	"testing"

	"github.com/opentasks/cmd/internal/query"
	"github.com/opentasks/cmd/internal/storage"
)

// E2ETestContext provides a fully-configured test environment for E2E tests
// Uses t.TempDir() for automatic cleanup
type E2ETestContext struct {
	Store   storage.BaseStorage
	Engine  *query.QueryEngine
	TmpDir  string
	Ctx     context.Context
	Cleanup func()
}

// SetupE2EEnvironment creates an isolated test environment with real file storage
// Follows the pattern from kubernetes/kubernetes and terraform for E2E CLI testing
func SetupE2EEnvironment(t *testing.T) *E2ETestContext {
	t.Helper()

	// Use t.TempDir() for automatic cleanup - no manual cleanup needed
	tmpDir := t.TempDir()

	// Create storage config pointing to temp directory
	storageConfig := storage.StorageConfig{
		Backend: "markdown-fs",
		Path:    tmpDir,
		Options: map[string]string{},
	}

	// Initialize storage backend
	store, err := storage.NewStorage(storageConfig)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Create query engine on top of storage
	engine := query.NewQueryEngine(store)

	// Create background context for operations
	ctx := context.Background()

	return &E2ETestContext{
		Store:  store,
		Engine: engine,
		TmpDir: tmpDir,
		Ctx:    ctx,
		Cleanup: func() {
			if err := store.Close(); err != nil {
				t.Logf("Warning: Failed to close storage: %v", err)
			}
		},
	}
}

// SetupMemoryEnvironment creates an isolated test environment with in-memory storage
// Useful for fast unit-style tests that don't need file system validation
func SetupMemoryEnvironment(t *testing.T) *E2ETestContext {
	t.Helper()

	// Create in-memory storage (no temp dir needed)
	store := storage.NewMemoryStorage()

	// Create query engine
	engine := query.NewQueryEngine(store)

	// Create background context
	ctx := context.Background()

	return &E2ETestContext{
		Store:  store,
		Engine: engine,
		TmpDir: "", // No temp dir for memory storage
		Ctx:    ctx,
		Cleanup: func() {
			if err := store.Close(); err != nil {
				t.Logf("Warning: Failed to close storage: %v", err)
			}
		},
	}
}
