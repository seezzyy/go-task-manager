package storage

import (
	"encoding/json"
	"os"
	"task-manager/internal/model"
)

type FileStorage struct {
	FilePath string
	Tasks    []model.Task
}

func (fs *FileStorage) save() error {
	data, err := json.MarshalIndent(fs.Tasks, "", "  ")

	if err != nil {
		return err
	}

	err = os.WriteFile(fs.FilePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FileStorage) load() error {
	data, err := os.ReadFile(fs.FilePath)

	if err != nil {
		if os.IsNotExist(err) {
			fs.Tasks = []model.Task{}
			return nil
		}
		return err
	}

	err = json.Unmarshal(data, &fs.Tasks)
	if err != nil {
		return err
	}

	return nil
}

func NewFileStorage(filePath string) (*FileStorage, error) {
	fs := &FileStorage{
		FilePath: filePath,
		Tasks:    []model.Task{},
	}

	err := fs.load()
	if err != nil {
		return nil, err
	}

	return fs, nil
}

func (fs *FileStorage) Add(task model.Task) (int, error) {
	maxID := 0

	for _, existingTask := range fs.Tasks {
		if existingTask.ID > maxID {
			maxID = existingTask.ID
		}
	}
	task.ID = maxID + 1

	fs.Tasks = append(fs.Tasks, task)

	err := fs.save()
	if err != nil {
		fs.Tasks = fs.Tasks[:len(fs.Tasks)-1]
		return 0, err

	}

	return task.ID, nil
}

func (fs *FileStorage) GetByID(id int) (model.Task, error) {
	for _, task := range fs.Tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return model.Task{}, model.ErrTaskNotFound
}

func (fs *FileStorage) Update(task model.Task) error {
	for i := range fs.Tasks {
		if fs.Tasks[i].ID == task.ID {
			oldTask := fs.Tasks[i]
			fs.Tasks[i] = task

			err := fs.save()

			if err != nil {
				fs.Tasks[i] = oldTask
				return err
			}

			return nil
		}
	}
	return model.ErrTaskNotFound
}

func (fs *FileStorage) Delete(id int) error {
	for i := range fs.Tasks {
		if fs.Tasks[i].ID == id {
			oldTasks := make([]model.Task, len(fs.Tasks))
			copy(oldTasks, fs.Tasks)
			fs.Tasks = append(
				fs.Tasks[:i],
				fs.Tasks[i+1:]...,
			)
			err := fs.save()
			if err != nil {
				fs.Tasks = oldTasks
				return err
			}
			return nil
		}
	}
	return model.ErrTaskNotFound
}

func (fs *FileStorage) List() ([]model.Task, error) {
	tasksCopy := make([]model.Task, len(fs.Tasks))
	copy(tasksCopy, fs.Tasks)

	return tasksCopy, nil
}
