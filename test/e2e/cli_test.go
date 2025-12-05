package e2e

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/opentasks/cmd/internal/model"
	"github.com/opentasks/cmd/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ExecuteCLI runs opentask CLI command in test environment and captures output
// Returns: stdout, stderr, exitCode
func ExecuteCLI(t *testing.T, tmpDir string, args ...string) (string, string, int) {
	t.Helper()

	// Find the opentask binary
	binaryPath := "opentask"
	binaryFound := false

	// Try multiple locations to find the binary
	possiblePaths := []string{
		"opentask",           // Try PATH first
		"./bin/opentask",     // Relative to current working directory
		"../bin/opentask",    // Up one directory
		"../../bin/opentask", // Up two directories
	}

	// Also try using runtime.Caller for module root
	if _, filename, _, ok := runtime.Caller(0); ok {
		moduleRoot := filepath.Join(filepath.Dir(filename), "..", "..")
		possiblePaths = append(possiblePaths, filepath.Join(moduleRoot, "bin", "opentask"))
	}

	// Find the first path that exists
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			// Convert to absolute path since we'll change working directories later
			absPath, err := filepath.Abs(path)
			if err != nil {
				// Fallback to original path if Abs fails
				binaryPath = path
			} else {
				binaryPath = absPath
			}
			binaryFound = true
			break
		}
	}

	// Skip the test if binary can't be found
	if !binaryFound {
		// Check if it's in PATH
		if _, err := exec.LookPath("opentask"); err != nil {
			t.Skipf("opentask binary not found in PATH or expected locations; skipping CLI test")
			return "", "", 0
		}
	}

	// Build command
	cmd := exec.Command(binaryPath, args...)

	// Capture output
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	// Set working directory to tmpDir
	cmd.Dir = tmpDir

	// Set environment to use tmpDir as project path
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("OPENTASK_PROJECT_PATH=%s", tmpDir),
	)

	// Execute
	exitCode := 0
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Logf("Command execution error: %v", err)
			exitCode = 1
		}
	}

	return outBuf.String(), errBuf.String(), exitCode
}

// TestE2E_CLITaskCreation tests task creation via CLI
func TestE2E_CLITaskCreation(t *testing.T) {
	env := SetupE2EEnvironment(t)
	defer env.Cleanup()

	if env.TmpDir == "" {
		t.Skip("CLI tests require file storage (TmpDir)")
	}

	// Initialize project config in tmpDir
	configPath := filepath.Join(env.TmpDir, ".opentask.toml")
	configContent := `[project]
id = "test-project"
name = "Test Project"
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	t.Run("create epic via CLI", func(t *testing.T) {
		stdout, stderr, exitCode := ExecuteCLI(t, env.TmpDir,
			"task", "new",
			"My Epic", // title is positional arg
			"--type", "epic",
			"--status", "planning",
		)

		assert.Equal(t, 0, exitCode, "CLI should exit 0")
		t.Logf("stdout: %s", stdout)
		t.Logf("stderr: %s", stderr)

		// Should mention task creation
		assert.Contains(t, strings.ToLower(stdout+stderr), "created",
			"Output should mention task creation")

		// Verify task exists in storage
		tasks, err := env.Store.ListTasks(env.Ctx)
		require.NoError(t, err)

		// Find the epic we just created
		var epic *model.Task
		for _, task := range tasks {
			if task.Title == "My Epic" {
				epic = task
				break
			}
		}

		require.NotNil(t, epic, "Epic should exist in storage")
		assert.Equal(t, "My Epic", epic.Title)
		assert.Equal(t, model.TypeEpic, epic.Type)
		assert.Equal(t, "planning", epic.Status)

		// Verify file created
		Assert(t, epic, env).FileExists()
	})

	t.Run("list tasks via CLI", func(t *testing.T) {
		// Create a few test tasks first
		testutil.CreateEpic(t, env.Store, testutil.WithTitle("Epic One"))
		testutil.CreateStory(t, env.Store, testutil.WithTitle("Story Two"))
		testutil.CreateTask(t, env.Store, testutil.WithTitle("Task Three"))

		stdout, stderr, exitCode := ExecuteCLI(t, env.TmpDir, "task", "list")

		assert.Equal(t, 0, exitCode, "CLI should exit 0")
		t.Logf("stdout:\n%s", stdout)

		// Output should contain task titles
		combinedOutput := stdout + stderr
		assert.Contains(t, combinedOutput, "Epic One")
		assert.Contains(t, combinedOutput, "Story Two")
		assert.Contains(t, combinedOutput, "Task Three")
	})

	t.Run("show task via CLI", func(t *testing.T) {
		// Create a task
		task := testutil.CreateTask(t, env.Store,
			testutil.WithTitle("Show Me Task"),
			testutil.WithDescription("This is a test description"),
		)

		stdout, stderr, exitCode := ExecuteCLI(t, env.TmpDir,
			"task", "show", fmt.Sprintf("%d", task.ID))

		assert.Equal(t, 0, exitCode, "CLI should exit 0")
		t.Logf("stdout:\n%s", stdout)

		combinedOutput := stdout + stderr
		assert.Contains(t, combinedOutput, "Show Me Task")
		assert.Contains(t, combinedOutput, "This is a test description")
	})

	t.Run("update task via CLI", func(t *testing.T) {
		// Create a task
		task := testutil.CreateTask(t, env.Store,
			testutil.WithTitle("Update Me"),
			testutil.WithStatus("todo"),
		)

		// Update status
		stdout, stderr, exitCode := ExecuteCLI(t, env.TmpDir,
			"task", "update",
			fmt.Sprintf("%d", task.ID),
			"--status", "done",
		)

		t.Logf("stdout: %s", stdout)
		t.Logf("stderr: %s", stderr)
		t.Logf("exitCode: %d", exitCode)

		// Verify in storage
		updated, err := env.Store.LoadTask(env.Ctx, task.ID)
		require.NoError(t, err)
		assert.Equal(t, "done", updated.Status)
	})
}

// TestE2E_CLIErrorHandling tests CLI error scenarios
func TestE2E_CLIErrorHandling(t *testing.T) {
	t.Run("no error when no project initialized (graceful)", func(t *testing.T) {
		// Create empty temp directory
		tmpDir := t.TempDir()

		stdout, stderr, exitCode := ExecuteCLI(t, tmpDir, "task", "list")

		// With active project resolution, operations should fail gracefully when no project is found
		// The CLI should exit with non-zero code but provide helpful error message
		assert.NotEqual(t, 0, exitCode, "CLI should exit with error when no project context")
		t.Logf("stdout: %s", stdout)
		t.Logf("stderr: %s", stderr)

		// Should show helpful error message about missing project
		combinedOutput := strings.ToLower(stdout + stderr)
		assert.True(t,
			strings.Contains(combinedOutput, "no active project") ||
				strings.Contains(combinedOutput, "not found") ||
				strings.Contains(combinedOutput, "project"),
			"Should show error message about missing project")
	})

	t.Run("error on invalid task ID", func(t *testing.T) {
		env := SetupE2EEnvironment(t)
		defer env.Cleanup()

		if env.TmpDir == "" {
			t.Skip("CLI tests require file storage")
		}

		// Initialize project
		configPath := filepath.Join(env.TmpDir, ".opentask.toml")
		err := os.WriteFile(configPath, []byte(`[project]
id = "test"
`), 0644)
		require.NoError(t, err)

		stdout, stderr, exitCode := ExecuteCLI(t, env.TmpDir,
			"task", "show", "999")

		assert.NotEqual(t, 0, exitCode, "Should exit non-zero for nonexistent task")
		t.Logf("stderr: %s", stderr)

		combinedOutput := strings.ToLower(stdout + stderr)
		assert.True(t,
			strings.Contains(combinedOutput, "not found") ||
				strings.Contains(combinedOutput, "error"),
			"Should show error message")
	})

	t.Run("error on missing required flags", func(t *testing.T) {
		env := SetupE2EEnvironment(t)
		defer env.Cleanup()

		if env.TmpDir == "" {
			t.Skip("CLI tests require file storage")
		}

		// Initialize project
		configPath := filepath.Join(env.TmpDir, ".opentask.toml")
		err := os.WriteFile(configPath, []byte(`[project]
id = "test"
`), 0644)
		require.NoError(t, err)

		// Try to create task without title
		stdout, stderr, exitCode := ExecuteCLI(t, env.TmpDir,
			"task", "new",
			"--type", "epic",
		)

		t.Logf("stdout: %s", stdout)
		t.Logf("stderr: %s", stderr)

		// Should show error (some CLIs exit non-zero, some show usage)
		combinedOutput := strings.ToLower(stdout + stderr)
		hasError := exitCode != 0 ||
			strings.Contains(combinedOutput, "required") ||
			strings.Contains(combinedOutput, "usage")

		assert.True(t, hasError, "Should indicate missing required flag")
	})

	t.Run("error on invalid task type", func(t *testing.T) {
		env := SetupE2EEnvironment(t)
		defer env.Cleanup()

		if env.TmpDir == "" {
			t.Skip("CLI tests require file storage")
		}

		// Initialize project
		configPath := filepath.Join(env.TmpDir, ".opentask.toml")
		err := os.WriteFile(configPath, []byte(`[project]
id = "test"
`), 0644)
		require.NoError(t, err)

		stdout, stderr, exitCode := ExecuteCLI(t, env.TmpDir,
			"task", "new",
			"Test",
			"--type", "invalid-type",
		)

		t.Logf("stdout: %s", stdout)
		t.Logf("stderr: %s", stderr)

		// Should show error
		combinedOutput := strings.ToLower(stdout + stderr)
		hasError := exitCode != 0 ||
			strings.Contains(combinedOutput, "invalid") ||
			strings.Contains(combinedOutput, "error")

		assert.True(t, hasError, "Should indicate invalid type")
	})
}
