package service

import (
	"strings"
	"task-manager/internal/model"
)

type TaskStorage interface {
	Add(task model.Task) (int, error)
	GetByID(id int) (model.Task, error)
	Update(task model.Task) error
	Delete(id int) error
	List() ([]model.Task, error)
}

type TaskService struct {
	Storage TaskStorage
}

func NewTaskService(storage TaskStorage) *TaskService {
	return &TaskService{
		Storage: storage,
	}
}

func (s *TaskService) CreateTask(
	title string,
	description string,
) (model.Task, error) {

	if err := validateTitle(title); err != nil {
		return model.Task{}, err
	}

	if err := validateDescription(description); err != nil {
		return model.Task{}, err
	}

	task := model.Task{
		Title:       title,
		Description: description,
		Status:      "new",
	}

	id, err := s.Storage.Add(task)
	if err != nil {
		return model.Task{}, err
	}

	task.ID = id
	return task, nil
}

func (s *TaskService) GetTask(id int) (model.Task, error) {
	task, err := s.Storage.GetByID(id)
	if err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (s *TaskService) CompleteTask(id int) error {
	task, err := s.Storage.GetByID(id)
	if err != nil {
		return err
	}
	task.Status = "completed"
	err = s.Storage.Update(task)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskService) DeleteTask(id int) error {
	return s.Storage.Delete(id)
}

func (s *TaskService) GetTasks() ([]model.Task, error) {
	tasks, err := s.Storage.List()
	if err != nil {
		return []model.Task{}, err
	}
	return tasks, nil
}

// UpdateTask partially updates a task, preserving fields
// that were not provided by the caller.
func (s *TaskService) UpdateTask(
	id int,
	title *string,
	description *string,
) error {
	if title == nil && description == nil {
		return model.ErrInvalidTask
	}

	task, err := s.Storage.GetByID(id)
	if err != nil {
		return err
	}

	if title != nil {
		if err := validateTitle(*title); err != nil {
			return err
		}

		task.Title = strings.TrimSpace(*title)
	}

	if description != nil {
		if err := validateDescription(*description); err != nil {
			return err
		}

		task.Description = strings.TrimSpace(*description)
	}
	return s.Storage.Update(task)

}
