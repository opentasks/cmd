package e2e

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opentasks/cmd/internal/model"
)

// AssertTask provides rich assertion failures for tasks
type AssertTask struct {
	t    *testing.T
	task *model.Task
	env  *E2ETestContext
}

// Assert creates a new task assertion helper
func Assert(t *testing.T, task *model.Task, env *E2ETestContext) *AssertTask {
	t.Helper()
	return &AssertTask{t: t, task: task, env: env}
}

// HasStatus asserts task status with rich failure output
func (a *AssertTask) HasStatus(expected string) *AssertTask {
	a.t.Helper()
	if a.task.Status != expected {
		a.t.Errorf(`
Task Status Mismatch
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Expected: %q
Actual:   %q

Task Details:
  ID:    %d
  Title: %q
  Type:  %q

Full Task:
%s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`,
			expected,
			a.task.Status,
			a.task.ID,
			a.task.Title,
			a.task.Type,
			a.formatTask())
	}
	return a
}

// HasTitle asserts task title
func (a *AssertTask) HasTitle(expected string) *AssertTask {
	a.t.Helper()
	if a.task.Title != expected {
		a.t.Errorf(`
Task Title Mismatch
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Expected: %q
Actual:   %q

Task: #%d (%s)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`,
			expected, a.task.Title, a.task.ID, a.task.Type)
	}
	return a
}

// HasParent asserts task has a parent relationship
func (a *AssertTask) HasParent(parentID int) *AssertTask {
	a.t.Helper()

	for _, rel := range a.task.Relationships {
		if rel.Type == model.RelParent && rel.TaskID == parentID {
			return a // Found it
		}
	}

	a.t.Errorf(`
Task Parent Relationship Missing
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Expected Parent ID: %d
Task: #%d %q (%s)

Current Relationships:
%s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`,
		parentID,
		a.task.ID,
		a.task.Title,
		a.task.Type,
		a.formatRelationships())
	return a
}

// HasChildCount asserts number of children via query engine
func (a *AssertTask) HasChildCount(expected int) *AssertTask {
	a.t.Helper()

	children, err := a.env.Engine.FindChildren(a.env.Ctx, a.task.ID)
	if err != nil {
		a.t.Fatalf("Failed to query children: %v", err)
	}

	actual := len(children)
	if actual != expected {
		childList := ""
		for i, child := range children {
			childList += fmt.Sprintf("  %d. #%d %q (%s) - %s\n",
				i+1, child.ID, child.Title, child.Type, child.Status)
		}

		a.t.Errorf(`
Child Count Mismatch
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Expected: %d children
Actual:   %d children

Parent Task: #%d %q (%s)

Children Found:
%s
Storage Path: %s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`,
			expected,
			actual,
			a.task.ID,
			a.task.Title,
			a.task.Type,
			childList,
			a.env.TmpDir)
	}
	return a
}

// HasType asserts task type with rich failure output
func (a *AssertTask) HasType(expected string) *AssertTask {
	a.t.Helper()
	if a.task.Type != expected {
		a.t.Errorf(`
Task Type Mismatch
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Expected: %q
Actual:   %q

Task: #%d %q
Status: %s

Valid types: epic, plan, research, story, decision, task
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`,
			expected, a.task.Type, a.task.ID, a.task.Title, a.task.Status)
	}
	return a
}

// HasDescription asserts task description content
func (a *AssertTask) HasDescription(expected string) *AssertTask {
	a.t.Helper()
	if a.task.Description != expected {
		a.t.Errorf(`
Task Description Mismatch
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Expected: %q
Actual:   %q

Task: #%d %q (%s)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`,
			expected, a.task.Description, a.task.ID, a.task.Title, a.task.Type)
	}
	return a
}

// HasTag asserts task has a specific tag
func (a *AssertTask) HasTag(tag string) *AssertTask {
	a.t.Helper()

	for _, t := range a.task.Tags {
		if t == tag {
			return a // Found it
		}
	}

	a.t.Errorf(`
Task Tag Missing
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Expected Tag: %q
Task: #%d %q (%s)

Current Tags: %v
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`,
		tag, a.task.ID, a.task.Title, a.task.Type, a.task.Tags)
	return a
}

// HasTags asserts task has all specified tags
func (a *AssertTask) HasTags(tags ...string) *AssertTask {
	a.t.Helper()

	missing := []string{}
	for _, expectedTag := range tags {
		found := false
		for _, actualTag := range a.task.Tags {
			if actualTag == expectedTag {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, expectedTag)
		}
	}

	if len(missing) > 0 {
		a.t.Errorf(`
Task Tags Missing
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Expected Tags: %v
Missing Tags:  %v
Task: #%d %q (%s)

Current Tags: %v
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`,
			tags, missing, a.task.ID, a.task.Title, a.task.Type, a.task.Tags)
	}
	return a
}

// HasBlocksRelationship asserts task has "blocks" relationship
func (a *AssertTask) HasBlocksRelationship(taskID int) *AssertTask {
	a.t.Helper()

	for _, rel := range a.task.Relationships {
		if rel.Type == model.RelBlocks && rel.TaskID == taskID {
			return a // Found it
		}
	}

	a.t.Errorf(`
Task "blocks" Relationship Missing
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Expected: blocks → Task #%d
Task: #%d %q (%s)

Current Relationships:
%s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`,
		taskID,
		a.task.ID,
		a.task.Title,
		a.task.Type,
		a.formatRelationships())
	return a
}

// HasRelatedToRelationship asserts task has "relates-to" relationship
func (a *AssertTask) HasRelatedToRelationship(taskID int) *AssertTask {
	a.t.Helper()

	for _, rel := range a.task.Relationships {
		if rel.Type == model.RelRelatedTo && rel.TaskID == taskID {
			return a // Found it
		}
	}

	a.t.Errorf(`
Task "relates-to" Relationship Missing
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Expected: relates-to → Task #%d
Task: #%d %q (%s)

Current Relationships:
%s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`,
		taskID,
		a.task.ID,
		a.task.Title,
		a.task.Type,
		a.formatRelationships())
	return a
}

// FileExists asserts task file exists on disk
func (a *AssertTask) FileExists() *AssertTask {
	a.t.Helper()

	if a.env.TmpDir == "" {
		// Memory storage, skip file check
		return a
	}

	// List all files in temp directory
	var files []string
	_ = filepath.Walk(a.env.TmpDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			relPath, _ := filepath.Rel(a.env.TmpDir, path)
			files = append(files, relPath)
		}
		return nil
	})

	// Check if file for this task exists
	// Files are named like: epic-1-auth-system.md or story-2-user-login.md
	found := false
	expectedFilePrefix := fmt.Sprintf("%s-%d-", a.task.Type, a.task.ID)
	for _, file := range files {
		if strings.HasPrefix(filepath.Base(file), expectedFilePrefix) && filepath.Ext(file) == ".md" {
			found = true
			break
		}
	}

	if !found {
		a.t.Errorf(`
Task File Not Found
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Task: #%d %q (%s)
Storage: %s

Files in storage:
  %v

Expected a markdown file to exist for this task.
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`,
			a.task.ID,
			a.task.Title,
			a.task.Type,
			a.env.TmpDir,
			files)
	}
	return a
}

// formatTask returns JSON representation of task
func (a *AssertTask) formatTask() string {
	data, _ := json.MarshalIndent(a.task, "  ", "  ")
	return string(data)
}

// formatRelationships returns human-readable relationship list
func (a *AssertTask) formatRelationships() string {
	if len(a.task.Relationships) == 0 {
		return "  (none)"
	}

	result := ""
	for i, rel := range a.task.Relationships {
		result += fmt.Sprintf("  %d. %s → Task #%d\n", i+1, rel.Type, rel.TaskID)
	}
	return result
}

// DumpEnvironment logs full test environment state for debugging
func DumpEnvironment(t *testing.T, env *E2ETestContext, msg string) {
	t.Helper()
	t.Logf(`
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Environment Dump: %s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Storage Path: %s

Files in storage:
%s

All tasks in storage:
%s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━`,
		msg,
		env.TmpDir,
		formatFileList(env.TmpDir),
		formatAllTasks(t, env))
}

func formatFileList(tmpDir string) string {
	if tmpDir == "" {
		return "  (memory storage - no files)"
	}

	var files []string
	_ = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			relPath, _ := filepath.Rel(tmpDir, path)
			size := info.Size()
			files = append(files, fmt.Sprintf("  - %s (%d bytes)", relPath, size))
		}
		return nil
	})

	if len(files) == 0 {
		return "  (no files)"
	}

	result := ""
	for _, f := range files {
		result += f + "\n"
	}
	return result
}

func formatAllTasks(t *testing.T, env *E2ETestContext) string {
	tasks, err := env.Store.ListTasks(env.Ctx)
	if err != nil {
		return fmt.Sprintf("  Error listing tasks: %v", err)
	}

	if len(tasks) == 0 {
		return "  (no tasks)"
	}

	result := ""
	for i, task := range tasks {
		result += fmt.Sprintf("  %d. #%d %q (%s) - %s\n",
			i+1, task.ID, task.Title, task.Type, task.Status)

		if len(task.Relationships) > 0 {
			for _, rel := range task.Relationships {
				result += fmt.Sprintf("     └─ %s: Task #%d\n", rel.Type, rel.TaskID)
			}
		}
	}
	return result
}
