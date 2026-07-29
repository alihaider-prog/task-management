package service

import (
	"errors"
	// "fmt"
	"task-management/internal/models"
	"task-management/internal/repository"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) GetTasks (ProjectID int64) ([]models.Task, error) {
	if ProjectID <= 0 {
		return nil, errors.New("invalid project id")
	}
	return s.repo.List(ProjectID)
}

func (s *TaskService) GetTaskByID(id int64) (*models.Task, error) {
	// implement optimistic locking using the version field.
	if id <= 0 {
		return nil, errors.New("invalid task id")
	}
	return s.repo.GetByID(id)
}

func (s *TaskService) CreateTask(task *models.Task) (*int64, error) {

	if task.Title == "" {
		return nil, errors.New("title is required")
	}

	if task.ProjectID == nil || *task.ProjectID <= 0 {
		return nil, errors.New("project is required")
	}

	if task.Priority == "" {
		task.Priority = models.PriorityMedium
	}
	
	task.Status = models.StatusTodo
	task.Version = 1
	task.CreatedBy = 1

	return s.repo.Create(task)
}

func (s *TaskService) UpdateTask(task *models.Task) error {
	if task.ID == 0 {
		return errors.New("invalid task")
	}
	return s.repo.Update(task)
}

func (s *TaskService) DeleteTask(id int64) error {
	if id <= 0 {
		return errors.New("invalid task id")
	}
	return s.repo.Delete(id)
}