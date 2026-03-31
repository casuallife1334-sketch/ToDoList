package service

import (
	"ToDoListNilchan/internal/core"
	"context"
	"fmt"
	"time"
)

type TaskRepository interface {
	Create(ctx context.Context, task core.TaskDomain) (core.TaskDomain, error)
	Get(ctx context.Context, id int) (core.TaskDomain, error)
	GetAll(ctx context.Context) ([]core.TaskDomain, error)
	Update(ctx context.Context, task core.TaskCompleteDomain) (core.TaskDomain, error)
	GetAllNotUpdated(ctx context.Context) ([]core.TaskDomain, error)
	Delete(ctx context.Context, id int) error
}

type Service struct {
	taskRepository TaskRepository
}

func NewService(taskRepository TaskRepository) *Service {
	return &Service{
		taskRepository: taskRepository,
	}
}

func (s *Service) CreateTask(ctx context.Context, title, description string) (core.TaskDomain, error) {
	task := core.TaskDomain{
		Title:       title,
		Description: description,
		Completed:   false,

		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}

	taskDomain, err := s.taskRepository.Create(ctx, task)
	if err != nil {
		return core.TaskDomain{}, fmt.Errorf("Create in repo: %w", err)
	}

	return taskDomain, nil
}

func (s *Service) GetTask(ctx context.Context, id int) (core.TaskDomain, error) {
	task, err := s.taskRepository.Get(ctx, id)
	if err != nil {
		return core.TaskDomain{}, fmt.Errorf("Get from repo: %w", err)
	}

	return task, nil
}

func (s *Service) GetAllTasks(ctx context.Context) ([]core.TaskDomain, error) {
	tasks, err := s.taskRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Service) GetAllUncompletedTasks(ctx context.Context) ([]core.TaskDomain, error) {
	tasks, err := s.taskRepository.GetAllNotUpdated(ctx)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Service) CompleteTask(ctx context.Context, id int) (core.TaskDomain, error) {
	initTime := time.Now()

	task := core.TaskCompleteDomain{
		ID:          id,
		Completed:   true,
		CompletedAt: &initTime,
	}

	updatedTask, err := s.taskRepository.Update(ctx, task)
	if err != nil {
		return core.TaskDomain{}, fmt.Errorf("Update in repo: %w", err)
	}

	return updatedTask, nil
}

func (s *Service) UncompleteTask(ctx context.Context, id int) (core.TaskDomain, error) {
	task := core.TaskCompleteDomain{
		ID:          id,
		Completed:   false,
		CompletedAt: nil,
	}

	updatedTask, err := s.taskRepository.Update(ctx, task)
	if err != nil {
		return core.TaskDomain{}, fmt.Errorf("Update in repo: %w", err)
	}

	return updatedTask, nil
}

func (s *Service) DeleteTask(ctx context.Context, id int) error {
	if err := s.taskRepository.Delete(ctx, id); err != nil {
		return fmt.Errorf("Delete in repo: %w", err)
	}

	return nil
}
