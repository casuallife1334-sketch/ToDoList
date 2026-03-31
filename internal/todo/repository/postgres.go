package repository

import (
	"ToDoListNilchan/internal/core"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) *Repository {
	return &Repository{
		conn: conn,
	}
}

func (r *Repository) Create(ctx context.Context, task core.TaskDomain) (core.TaskDomain, error) {
	sqlQuery := `
	INSERT INTO tasks(title, description, completed, created_at, completed_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, title, description, completed, created_at, completed_at
	`

	var taskDomain core.TaskDomain

	if err := r.conn.QueryRow(
		ctx,
		sqlQuery,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
	).Scan(
		&taskDomain.ID,
		&taskDomain.Title,
		&taskDomain.Description,
		&taskDomain.Completed,
		&taskDomain.CreatedAt,
		&taskDomain.CompletedAt,
	); err != nil {
		return core.TaskDomain{}, core.ErrNotFound
	}

	return taskDomain, nil
}

func (r *Repository) Get(ctx context.Context, id int) (core.TaskDomain, error) {
	sqlQuery := `
	SELECT id, title, description, completed, created_at, completed_at FROM tasks
	WHERE id=$1
	`

	var task core.TaskDomain

	if err := r.conn.QueryRow(
		ctx,
		sqlQuery,
		id,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.CompletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.TaskDomain{}, core.ErrNotFound
		} else {
			return core.TaskDomain{}, err
		}
	}

	return task, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]core.TaskDomain, error) {
	sqlQuery := `
	SELECT id, title, description, completed, created_at, completed_at FROM tasks
	`

	rows, err := r.conn.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]core.TaskDomain, 0)

	for rows.Next() {
		var task core.TaskDomain

		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.CompletedAt,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *Repository) GetAllNotUpdated(ctx context.Context) ([]core.TaskDomain, error) {
	sqlQuery := `
	SELECT id, title, description, completed, created_at, completed_at FROM tasks
	WHERE NOT completed
	`

	rows, err := r.conn.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notUpdatedTasks []core.TaskDomain

	for rows.Next() {
		var notUpdatedTask core.TaskDomain

		if err := rows.Scan(
			&notUpdatedTask.ID,
			&notUpdatedTask.Title,
			&notUpdatedTask.Description,
			&notUpdatedTask.Completed,
			&notUpdatedTask.CreatedAt,
			&notUpdatedTask.CompletedAt,
		); err != nil {
			return nil, err
		}

		notUpdatedTasks = append(notUpdatedTasks, notUpdatedTask)
	}

	return notUpdatedTasks, nil
}

func (r *Repository) Update(ctx context.Context, task core.TaskCompleteDomain) (core.TaskDomain, error) {
	sqlQuery := `
	UPDATE tasks
	SET completed=$2, completed_at=$3
	WHERE id=$1
	RETURNING id, title, description, completed, created_at, completed_at
	`

	var taskDomain core.TaskDomain

	if err := r.conn.QueryRow(
		ctx,
		sqlQuery,
		task.ID,
		task.Completed,
		task.CompletedAt,
	).Scan(
		&taskDomain.ID,
		&taskDomain.Title,
		&taskDomain.Description,
		&taskDomain.Completed,
		&taskDomain.CreatedAt,
		&taskDomain.CompletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.TaskDomain{}, core.ErrNotFound
		} else {
			return core.TaskDomain{}, err
		}
	}

	return taskDomain, nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	sqlQuery := `
	DELETE FROM tasks
	WHERE id=$1
	`

	if _, err := r.conn.Exec(ctx, sqlQuery, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.ErrNotFound
		} else {
			return err
		}
	}

	return nil
}
