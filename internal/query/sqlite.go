package query

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/opentasks/cmd/internal/model"
	"github.com/opentasks/cmd/internal/storage"
)

// SQLiteService implements storage.BaseStorage using an in-memory SQLite database
// It loads markdown files into SQLite and normalizes relationships into join tables,
// enabling DAG queries with SQL syntax.
type SQLiteService struct {
	db            *sql.DB
	sourceStorage storage.BaseStorage // Underlying markdown file storage for reads/writes
}

// NewSQLiteService creates a new SQLite-backed storage service
// It loads all tasks from sourceStorage into memory and builds relationship indices
func NewSQLiteService(ctx context.Context, sourceStorage storage.BaseStorage) (*SQLiteService, error) {
	// Create in-memory SQLite database
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		return nil, fmt.Errorf("failed to create in-memory SQLite: %w", err)
	}

	// Create schema
	if err := createSchema(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	svc := &SQLiteService{
		db:            db,
		sourceStorage: sourceStorage,
	}

	// Load all tasks from source storage into SQLite
	if err := svc.loadAllTasks(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}

	return svc, nil
}

// createSchema creates the SQLite schema for tasks and relationships
func createSchema(ctx context.Context, db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		type TEXT NOT NULL,
		status TEXT NOT NULL,
		description TEXT,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS task_tags (
		task_id INTEGER NOT NULL,
		tag TEXT NOT NULL,
		FOREIGN KEY (task_id) REFERENCES tasks(id),
		PRIMARY KEY (task_id, tag)
	);

	CREATE TABLE IF NOT EXISTS task_relationships (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		from_task_id INTEGER NOT NULL,
		to_task_id INTEGER NOT NULL,
		relation_type TEXT NOT NULL,
		FOREIGN KEY (from_task_id) REFERENCES tasks(id),
		FOREIGN KEY (to_task_id) REFERENCES tasks(id),
		UNIQUE (from_task_id, to_task_id, relation_type)
	);

	CREATE INDEX IF NOT EXISTS idx_relationships_from ON task_relationships(from_task_id);
	CREATE INDEX IF NOT EXISTS idx_relationships_to ON task_relationships(to_task_id);
	CREATE INDEX IF NOT EXISTS idx_relationships_type ON task_relationships(relation_type);
	CREATE INDEX IF NOT EXISTS idx_tags_task ON task_tags(task_id);
	`

	_, err := db.ExecContext(ctx, schema)
	return err
}

// loadAllTasks loads all tasks from source storage into SQLite
func (s *SQLiteService) loadAllTasks(ctx context.Context) error {
	tasks, err := s.sourceStorage.ListTasks(ctx)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if err := s.insertTask(ctx, task); err != nil {
			return fmt.Errorf("failed to insert task %d: %w", task.ID, err)
		}
	}

	return nil
}

// insertTask inserts a task and its relationships/tags into SQLite
func (s *SQLiteService) insertTask(ctx context.Context, task *model.Task) error {
	// Insert task
	_, err := s.db.ExecContext(ctx,
		`INSERT OR REPLACE INTO tasks (id, title, type, status, description, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		task.ID, task.Title, task.Type, task.Status, task.Description,
		task.CreatedAt.Format(time.RFC3339), task.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return err
	}

	// Insert tags
	for _, tag := range task.Tags {
		_, err := s.db.ExecContext(ctx,
			`INSERT OR IGNORE INTO task_tags (task_id, tag) VALUES (?, ?)`,
			task.ID, tag,
		)
		if err != nil {
			return err
		}
	}

	// Insert relationships as join table entries
	for _, rel := range task.Relationships {
		_, err := s.db.ExecContext(ctx,
			`INSERT OR REPLACE INTO task_relationships (from_task_id, to_task_id, relation_type)
			 VALUES (?, ?, ?)`,
			task.ID, rel.TaskID, rel.Type,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadTask retrieves a task by ID from SQLite
func (s *SQLiteService) LoadTask(ctx context.Context, id int) (*model.Task, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, title, type, status, description, created_at, updated_at FROM tasks WHERE id = ?`,
		id,
	)

	task := &model.Task{}
	createdAt, updatedAt := "", ""

	err := row.Scan(&task.ID, &task.Title, &task.Type, &task.Status, &task.Description, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, storage.ErrTaskNotFound
	}
	if err != nil {
		return nil, err
	}

	// Parse timestamps
	task.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	task.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	// Load tags
	tagRows, err := s.db.QueryContext(ctx, `SELECT tag FROM task_tags WHERE task_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer tagRows.Close()

	for tagRows.Next() {
		var tag string
		if err := tagRows.Scan(&tag); err != nil {
			return nil, err
		}
		task.Tags = append(task.Tags, tag)
	}

	// Load relationships
	relRows, err := s.db.QueryContext(ctx,
		`SELECT to_task_id, relation_type FROM task_relationships WHERE from_task_id = ?`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer relRows.Close()

	for relRows.Next() {
		rel := model.Relationship{}
		if err := relRows.Scan(&rel.TaskID, &rel.Type); err != nil {
			return nil, err
		}
		task.Relationships = append(task.Relationships, rel)
	}

	return task, nil
}

// SaveTask persists a task back to source storage and updates SQLite cache
func (s *SQLiteService) SaveTask(ctx context.Context, task *model.Task) error {
	// Save to source storage (markdown files)
	if err := s.sourceStorage.SaveTask(ctx, task); err != nil {
		return err
	}

	// Update SQLite cache
	return s.insertTask(ctx, task)
}

// DeleteTask removes a task from both source storage and SQLite cache
func (s *SQLiteService) DeleteTask(ctx context.Context, id int) error {
	// Delete from source storage
	if err := s.sourceStorage.DeleteTask(ctx, id); err != nil {
		return err
	}

	// Delete from SQLite
	_, err := s.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `DELETE FROM task_tags WHERE task_id = ?`, id)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `DELETE FROM task_relationships WHERE from_task_id = ? OR to_task_id = ?`, id, id)
	return err
}

// ListTasks returns tasks matching filters
// Falls back to source storage filters for now, could be enhanced with SQL WHERE clauses
func (s *SQLiteService) ListTasks(ctx context.Context, filters ...storage.TaskFilter) ([]*model.Task, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, title, type, status, description, created_at, updated_at FROM tasks ORDER BY created_at`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		task := &model.Task{}
		createdAt, updatedAt := "", ""

		if err := rows.Scan(&task.ID, &task.Title, &task.Type, &task.Status, &task.Description, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		task.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		task.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		// Load tags
		tagRows, err := s.db.QueryContext(ctx, `SELECT tag FROM task_tags WHERE task_id = ?`, task.ID)
		if err != nil {
			return nil, err
		}
		for tagRows.Next() {
			var tag string
			if err := tagRows.Scan(&tag); err != nil {
				_ = tagRows.Close()
				return nil, err
			}
			task.Tags = append(task.Tags, tag)
		}
		_ = tagRows.Close()

		// Load relationships
		relRows, err := s.db.QueryContext(ctx,
			`SELECT to_task_id, relation_type FROM task_relationships WHERE from_task_id = ?`,
			task.ID,
		)
		if err != nil {
			return nil, err
		}
		for relRows.Next() {
			rel := model.Relationship{}
			if err := relRows.Scan(&rel.TaskID, &rel.Type); err != nil {
				_ = relRows.Close()
				return nil, err
			}
			task.Relationships = append(task.Relationships, rel)
		}
		_ = relRows.Close()

		tasks = append(tasks, task)
	}

	// Apply filters (for backward compatibility with existing filter functions)
	var filtered []*model.Task
	for _, task := range tasks {
		match := true
		for _, filter := range filters {
			if !filter(task) {
				match = false
				break
			}
		}
		if match {
			filtered = append(filtered, task)
		}
	}

	return filtered, nil
}

// NextID generates the next global sequential ID
func (s *SQLiteService) NextID(ctx context.Context) (int, error) {
	var maxID int
	err := s.db.QueryRowContext(ctx, `SELECT COALESCE(MAX(id), 0) FROM tasks`).Scan(&maxID)
	if err != nil {
		return 0, err
	}
	return maxID + 1, nil
}

// GetRelated returns all tasks related to the given task by relationship type
func (s *SQLiteService) GetRelated(ctx context.Context, taskID int, relationType string) ([]*model.Task, error) {
	// Find all tasks that have a relationship TO taskID of the given type
	rows, err := s.db.QueryContext(ctx,
		`SELECT DISTINCT t.id FROM tasks t
		 JOIN task_relationships tr ON t.id = tr.from_task_id
		 WHERE tr.to_task_id = ? AND tr.relation_type = ?
		 ORDER BY t.created_at`,
		taskID, relationType,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var related []*model.Task
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		task, err := s.LoadTask(ctx, id)
		if err != nil {
			return nil, err
		}
		related = append(related, task)
	}

	return related, nil
}

// QuerySQL executes a raw SQL query and returns task IDs
// This enables DAG queries like finding descendants with recursive CTEs
func (s *SQLiteService) QuerySQL(ctx context.Context, query string, args ...interface{}) ([]int, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// QueryDAG finds all descendants of a task (transitive parent relationships)
func (s *SQLiteService) QueryDAG(ctx context.Context, taskID int) ([]*model.Task, error) {
	// Use recursive CTE to find all descendants
	query := `
	WITH RECURSIVE descendants AS (
		SELECT id FROM tasks WHERE id IN (
			SELECT from_task_id FROM task_relationships
			WHERE to_task_id = ? AND relation_type = 'parent'
		)
		UNION ALL
		SELECT t.id FROM tasks t
		JOIN descendants d ON t.id IN (
			SELECT from_task_id FROM task_relationships
			WHERE to_task_id = d.id AND relation_type = 'parent'
		)
	)
	SELECT id FROM descendants
	`

	ids, err := s.QuerySQL(ctx, query, taskID)
	if err != nil {
		return nil, err
	}

	var descendants []*model.Task
	for _, id := range ids {
		task, err := s.LoadTask(ctx, id)
		if err == nil {
			descendants = append(descendants, task)
		}
	}

	return descendants, nil
}

// Close closes the SQLite connection
func (s *SQLiteService) Close() error {
	return s.db.Close()
}
