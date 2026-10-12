package storage

import (
	"task-manager/internal/model"
)

type MemoryStorage struct {
	Tasks []model.Task
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		Tasks: []model.Task{},
	}
}

func (ms *MemoryStorage) Add(task model.Task) (int, error) {
	maxID := 0

	for _, existingTask := range ms.Tasks {
		if existingTask.ID > maxID {
			maxID = existingTask.ID
		}
	}
	task.ID = maxID + 1
	ms.Tasks = append(ms.Tasks, task)
	return task.ID, nil
}

func (ms *MemoryStorage) GetByID(id int) (model.Task, error) {
	for _, task := range ms.Tasks {
		if task.ID == id {
			return task, nil

		}
	}
	return model.Task{}, model.ErrTaskNotFound
}

func (ms *MemoryStorage) Update(task model.Task) error {
	for i := range ms.Tasks {
		if ms.Tasks[i].ID == task.ID {
			ms.Tasks[i] = task
			return nil
		}
	}
	return model.ErrTaskNotFound
}

func (ms *MemoryStorage) Delete(id int) error {
	for i := range ms.Tasks {
		if ms.Tasks[i].ID == id {
			ms.Tasks = append(
				ms.Tasks[:i],
				ms.Tasks[i+1:]...,
			)
			return nil
		}
	}
	return model.ErrTaskNotFound
}

func (ms *MemoryStorage) List() ([]model.Task, error) {
	tasksCopy := make([]model.Task, len(ms.Tasks))
	copy(tasksCopy, ms.Tasks)

	return tasksCopy, nil
}
