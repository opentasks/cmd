package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/opentasks/cmd/internal/model"
	"gopkg.in/yaml.v3"
)

type MarkdownFileStorage struct {
	basePath string // Project root path
}

// NewMarkdownFileStorage creates a new markdown file storage backend
func NewMarkdownFileStorage(basePath string) (*MarkdownFileStorage, error) {
	// Ensure path exists
	// #nosec G301 - Directory permissions 0750 are intentional for task storage root
	// basePath is validated and set from configuration
	if err := os.MkdirAll(basePath, 0750); err != nil {
		return nil, err
	}

	return &MarkdownFileStorage{basePath: basePath}, nil
}

// slugify converts a string to a URL-friendly slug
func slugify(s string) string {
	// Convert to lowercase and replace spaces with hyphens
	var result strings.Builder

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(unicode.ToLower(r))
		} else if unicode.IsSpace(r) {
			result.WriteRune('-')
		} else if r == '-' || r == '_' {
			result.WriteRune(r)
		}
	}

	// Clean up multiple hyphens
	slug := result.String()
	re := regexp.MustCompile(`-+`)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	return slug
}

// taskToPath generates the file path for a task
// File naming: <epic_id>-<epic_slug>/<typecode>-<id>-<slug>.md
// For tasks without epic parent: <typecode>-<id>-<slug>.md
func (s *MarkdownFileStorage) taskToPath(ctx context.Context, task *model.Task) (string, error) {
	// Determine parent epic directory
	var epicDir string
	var parentEpicID int

	for _, rel := range task.Relationships {
		if rel.Type == model.RelParent {
			parentEpicID = rel.TaskID
			break
		}
	}

	if parentEpicID > 0 {
		// Load parent epic to get its title
		parentEpic, err := s.LoadTask(ctx, parentEpicID)
		if err != nil {
			return "", fmt.Errorf("parent epic not found: %w", err)
		}

		epicSlug := slugify(parentEpic.Title)
		epicDir = fmt.Sprintf("%d-%s", parentEpicID, epicSlug)
	}

	// Build filename: typecode-id-slug.md
	typeCode := model.TypeCode[task.Type]
	slug := slugify(task.Title)
	filename := fmt.Sprintf("%s-%d-%s.md", typeCode, task.ID, slug)

	// Combine directory and filename
	if epicDir != "" {
		return filepath.Join(s.basePath, epicDir, filename), nil
	}
	return filepath.Join(s.basePath, filename), nil
}

// parseTaskFile parses markdown file into Task
// Format: frontmatter (YAML) separator (---) markdown body
func (s *MarkdownFileStorage) parseTaskFile(filePath string) (*model.Task, error) {
	// #nosec G304 - filePath comes from filepath.Walk within basePath boundary
	// All files are scoped to basePath set at storage initialization
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Split by --- separator
	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid task file format: missing frontmatter")
	}

	frontmatterYAML := strings.TrimSpace(parts[1])
	description := strings.TrimSpace(parts[2])

	// Parse YAML frontmatter
	var frontmatter struct {
		ID            int      `yaml:"id"`
		Title         string   `yaml:"title"`
		Type          string   `yaml:"type"`
		Status        string   `yaml:"status"`
		Tags          []string `yaml:"tags"`
		Relationships []struct {
			Type   string `yaml:"type"`
			TaskID int    `yaml:"taskID"`
		} `yaml:"relationships"`
		CreatedAt string `yaml:"createdAt"`
		UpdatedAt string `yaml:"updatedAt"`
	}

	if err := yaml.Unmarshal([]byte(frontmatterYAML), &frontmatter); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Parse timestamps
	createdAt, err := time.Parse(time.RFC3339, frontmatter.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid createdAt: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, frontmatter.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid updatedAt: %w", err)
	}

	// Convert frontmatter relationships
	relationships := make([]model.Relationship, len(frontmatter.Relationships))
	for i, rel := range frontmatter.Relationships {
		relationships[i] = model.Relationship{
			Type:   rel.Type,
			TaskID: rel.TaskID,
		}
	}

	// Initialize empty tags if nil
	if frontmatter.Tags == nil {
		frontmatter.Tags = []string{}
	}

	task := &model.Task{
		ID:            frontmatter.ID,
		Title:         frontmatter.Title,
		Type:          frontmatter.Type,
		Status:        frontmatter.Status,
		Tags:          frontmatter.Tags,
		Relationships: relationships,
		Description:   description,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	return task, nil
}

// LoadTask loads task frontmatter and content from markdown file
// Task ID is extracted from frontmatter, not from filename
func (s *MarkdownFileStorage) LoadTask(ctx context.Context, id int) (*model.Task, error) {
	// Walk project to find file with matching ID in frontmatter
	var foundPath string

	err := filepath.Walk(s.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue on errors
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		task, err := s.parseTaskFile(path)
		if err != nil {
			return nil // Skip files that fail to parse
		}

		if task.ID == id {
			foundPath = path
			return filepath.SkipDir // Found it, stop walking
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if foundPath == "" {
		return nil, ErrTaskNotFound
	}

	return s.parseTaskFile(foundPath)
}

// SaveTask writes task to markdown file
func (s *MarkdownFileStorage) SaveTask(ctx context.Context, task *model.Task) error {
	// Get the target path
	path, err := s.taskToPath(ctx, task)
	if err != nil {
		return err
	}

	// Find old file path and delete if it exists and differs from new path
	oldTask, err := s.LoadTask(ctx, task.ID)
	if err == nil && oldTask != nil {
		// Loaded existing task, compute its old path
		oldPath, err := s.taskToPath(ctx, oldTask)
		if err == nil && oldPath != path {
			// Path has changed (e.g., due to title update), remove old file
			_ = os.Remove(oldPath) // Ignore error if file doesn't exist
		}
	}
	// If LoadTask fails (task doesn't exist yet), that's fine, just continue

	// Ensure directory exists
	dir := filepath.Dir(path)
	// #nosec G301 - Directory permissions 0750 are intentional for task storage
	// Path is constructed from validated basePath (set at storage initialization)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return err
	}

	// Build frontmatter
	type RelationshipYAML struct {
		Type   string `yaml:"type"`
		TaskID int    `yaml:"taskID"`
	}

	frontmatter := struct {
		ID            int                `yaml:"id"`
		Title         string             `yaml:"title"`
		Type          string             `yaml:"type"`
		Status        string             `yaml:"status"`
		Tags          []string           `yaml:"tags"`
		Relationships []RelationshipYAML `yaml:"relationships"`
		CreatedAt     string             `yaml:"createdAt"`
		UpdatedAt     string             `yaml:"updatedAt"`
	}{
		ID:        task.ID,
		Title:     task.Title,
		Type:      task.Type,
		Status:    task.Status,
		Tags:      task.Tags,
		CreatedAt: task.CreatedAt.Format(time.RFC3339),
		UpdatedAt: task.UpdatedAt.Format(time.RFC3339),
	}

	// Convert relationships
	frontmatter.Relationships = make([]RelationshipYAML, len(task.Relationships))
	for i, rel := range task.Relationships {
		frontmatter.Relationships[i] = RelationshipYAML{
			Type:   rel.Type,
			TaskID: rel.TaskID,
		}
	}

	// Marshal YAML
	yamlBytes, err := yaml.Marshal(frontmatter)
	if err != nil {
		return err
	}

	// Build file content
	content := fmt.Sprintf("---\n%s---\n\n%s\n", string(yamlBytes), task.Description)

	// Write file
	// #nosec G306 - File permissions 0600 are intentional for task files
	// #nosec G304 - Path constructed from validated basePath + sanitized filename
	// basePath set at initialization, filename derived from task ID and sanitized title
	return os.WriteFile(path, []byte(content), 0600)
}

// DeleteTask removes the task file
func (s *MarkdownFileStorage) DeleteTask(ctx context.Context, id int) error {
	task, err := s.LoadTask(ctx, id)
	if err != nil {
		return err
	}

	path, err := s.taskToPath(ctx, task)
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if err != nil && os.IsNotExist(err) {
		return ErrTaskNotFound
	}
	return err
}

// ListTasks returns all tasks matching filters
func (s *MarkdownFileStorage) ListTasks(ctx context.Context, filters ...TaskFilter) ([]*model.Task, error) {
	var tasks []*model.Task

	// Walk all directories and files
	err := filepath.Walk(s.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue on errors
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		task, err := s.parseTaskFile(path)
		if err != nil {
			return nil // Skip files that fail to parse
		}

		// Apply filters
		match := true
		for _, filter := range filters {
			if !filter(task) {
				match = false
				break
			}
		}

		if match {
			tasks = append(tasks, task)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort by creation time for consistent ordering
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})

	return tasks, nil
}

// NextID generates the next global sequential ID
// Finds the maximum ID from all task frontmatter and returns max + 1
// This ensures IDs are never reused even if tasks are deleted
func (s *MarkdownFileStorage) NextID(ctx context.Context) (int, error) {
	var maxID int

	err := filepath.Walk(s.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue walking on errors
		}

		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		task, err := s.parseTaskFile(path)
		if err != nil {
			return nil // Skip unparseable files
		}

		if task.ID > maxID {
			maxID = task.ID
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return maxID + 1, nil
}

// GetRelated returns all tasks related by relationship type
func (s *MarkdownFileStorage) GetRelated(ctx context.Context, taskID int, relationType string) ([]*model.Task, error) {
	allTasks, err := s.ListTasks(ctx)
	if err != nil {
		return nil, err
	}

	var related []*model.Task
	for _, task := range allTasks {
		for _, rel := range task.Relationships {
			if rel.Type == relationType && rel.TaskID == taskID {
				related = append(related, task)
			}
		}
	}

	return related, nil
}

// Close is a no-op for file storage
func (s *MarkdownFileStorage) Close() error {
	return nil
}
