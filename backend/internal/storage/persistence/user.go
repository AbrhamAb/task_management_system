package persistence

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"task-management-backend/internal/model/db"
)

func (s *Store) CreateUser(ctx context.Context, name, email, passwordHash string) (db.User, error) {
	q := `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, password_hash, created_at, updated_at
	`

	var user db.User
	err := s.db.QueryRowContext(ctx, q, name, email, passwordHash).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	q := `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user db.User
	err := s.db.QueryRowContext(ctx, q, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}

func (s *Store) GetUserByID(ctx context.Context, userID int64) (db.User, error) {
	q := `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user db.User
	err := s.db.QueryRowContext(ctx, q, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}

func (s *Store) ListTasksByUserID(ctx context.Context, userID int64) ([]db.Task, error) {
	q := `
		SELECT id, user_id, title, COALESCE(description, ''), status, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`

	rows, err := s.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]db.Task, 0)
	for rows.Next() {
		var task db.Task
		if err = rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (s *Store) CreateTask(ctx context.Context, userID int64, title, description, status string) (db.Task, error) {
	q := `
		INSERT INTO tasks (user_id, title, description, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, title, COALESCE(description, ''), status, created_at, updated_at
	`

	var task db.Task
	err := s.db.QueryRowContext(ctx, q, userID, title, description, status).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	return task, err
}

func (s *Store) UpdateTask(ctx context.Context, userID, taskID int64, title, description, status string) (db.Task, error) {
	q := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, updated_at = NOW()
		WHERE id = $4 AND user_id = $5
		RETURNING id, user_id, title, COALESCE(description, ''), status, created_at, updated_at
	`

	var task db.Task
	err := s.db.QueryRowContext(ctx, q, title, description, status, taskID, userID).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	return task, err
}

func (s *Store) UpdateTaskStatus(ctx context.Context, userID, taskID int64, status string) (db.Task, error) {
	q := `
		UPDATE tasks
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
		RETURNING id, user_id, title, COALESCE(description, ''), status, created_at, updated_at
	`

	var task db.Task
	err := s.db.QueryRowContext(ctx, q, status, taskID, userID).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	return task, err
}

func (s *Store) DeleteTask(ctx context.Context, userID, taskID int64) error {
	q := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`
	result, err := s.db.ExecContext(ctx, q, taskID, userID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "duplicate key value")
}

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
