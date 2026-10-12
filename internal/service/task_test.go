package service_test

import (
	"errors"
	"strings"
	"testing"

	"task-manager/internal/model"
	"task-manager/internal/service"
	"task-manager/internal/storage"
)

func TestCreateTask_EmptyTitle(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	taskService := service.NewTaskService(memoryStorage)

	_, err := taskService.CreateTask("", "Some description")

	if !errors.Is(err, model.ErrInvalidTask) {
		t.Fatalf(
			"expected ErrInvalidTask, got %v",
			err,
		)
	}
}

func TestCreateTask_Success(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	taskService := service.NewTaskService(memoryStorage)

	task, err := taskService.CreateTask(
		"Learn Go",
		"Study unit testing",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}

	if task.Title != "Learn Go" {
		t.Errorf("expected Title Learn Go, got %q", task.Title)
	}

	if task.Description != "Study unit testing" {
		t.Errorf("expected Description 'Study unit testing', got %q", task.Description)
	}

	if task.Status != "new" {
		t.Errorf("expected Status, got %q", task.Status)
	}
}

func TestGetTask_Success(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	taskService := service.NewTaskService(memoryStorage)

	createTask, err := taskService.CreateTask(
		"Learn Go",
		"Testing GetTask",
	)

	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	receivedTask, err := taskService.GetTask(createTask.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedTask.ID != createTask.ID {
		t.Errorf("expected ID %d, got %d", createTask.ID, receivedTask.ID)
	}

	if receivedTask.Title != "Learn Go" {
		t.Errorf("expected Title 'Learn Go', got %q", receivedTask.Title)
	}
}

func TestDeleteTask_Success(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	taskService := service.NewTaskService(memoryStorage)

	createTask, err := taskService.CreateTask(
		"Learn Go",
		"Testing DeleteTask",
	)

	if err != nil {
		t.Fatalf("failed to create: %v", err)
	}

	err = taskService.DeleteTask(createTask.ID)

	if err != nil {
		t.Fatalf("failed to delete task: %v", err)
	}

	_, err = taskService.GetTask(createTask.ID)

	if !errors.Is(err, model.ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	taskService := service.NewTaskService(memoryStorage)

	_, err := taskService.GetTask(999)

	if !errors.Is(err, model.ErrTaskNotFound) {
		t.Fatalf(
			"expected ErrTaskNotFound, got %v",
			err,
		)
	}
}
func TestDeleteTask_NotFound(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	taskService := service.NewTaskService(memoryStorage)

	err := taskService.DeleteTask(999)
	if !errors.Is(err, model.ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestUpdateTask_Success(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	taskService := service.NewTaskService(memoryStorage)

	createdTask, err := taskService.CreateTask(
		"Learn Go",
		"Testing UpdateTask",
	)

	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	newTitle := "Learn Advanced Go"

	err = taskService.UpdateTask(
		createdTask.ID,
		&newTitle,
		nil,
	)

	if err != nil {
		t.Fatalf("failed to update task: %v", err)
	}

	updatedTask, err := taskService.GetTask(createdTask.ID)

	if err != nil {
		t.Fatalf("failed to get update task: %v", err)
	}

	if updatedTask.Title != newTitle {
		t.Errorf("expected title %q, got %q", newTitle, updatedTask.Title)
	}

	if updatedTask.Description != "Testing UpdateTask" {
		t.Errorf("expected descriptions %q, got %q", "Testing UpdateTask", updatedTask.Description)
	}
}

func TestUpdateTask_DescriptionOnly(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	taskService := service.NewTaskService(memoryStorage)

	createdTask, err := taskService.CreateTask(
		"Learn Go",
		"Old Description",
	)

	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	newDescription := "New Description"

	err = taskService.UpdateTask(
		createdTask.ID,
		nil,
		&newDescription,
	)

	if err != nil {
		t.Fatalf("failed to update tasl: %v", err)
	}

	updatedTask, err := taskService.GetTask(createdTask.ID)

	if err != nil {
		t.Fatalf("failed to get update task: %v", err)
	}

	if updatedTask.Description != newDescription {
		t.Errorf("expected description %q, got %q", newDescription, updatedTask.Description)
	}

	if updatedTask.Title != "Learn Go" {
		t.Errorf("expected title %q, got %q", "Learn Go", updatedTask.Title)
	}
}

func TestUpdateTask_EmptyTitle(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	taskService := service.NewTaskService(memoryStorage)

	createdTask, err := taskService.CreateTask(
		"Learn Go",
		"Study testing",
	)

	if err != nil {
		t.Fatalf("failed to created task: %v", err)
	}

	emptyTitle := ""

	err = taskService.UpdateTask(
		createdTask.ID,
		&emptyTitle,
		nil,
	)

	if !errors.Is(err, model.ErrInvalidTask) {
		t.Fatalf("expected ErrInvalidTask, got %v", err)
	}

	nameTask, err := taskService.GetTask(createdTask.ID)

	if err != nil {
		t.Fatalf("failed to update tsk: %v", err)
	}

	if nameTask.Title != "Learn Go" {
		t.Errorf("expected title %q, got %q", "Learn Go", nameTask.Title)
	}

}

func TestCreateTask_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{
			name:    "empty title",
			title:   "",
			wantErr: true,
		},
		{
			name:    "space only",
			title:   "    ",
			wantErr: true,
		},
		{
			name:    "valid title",
			title:   "Learn Go",
			wantErr: false,
		},
		{
			name:    "tabs only",
			title:   "\t\t",
			wantErr: true,
		},
		{
			name:    "newline only",
			title:   "\n",
			wantErr: true,
		},
		{
			name:    "valid title with spaces",
			title:   "Learn Go Backend",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memoryStorage := storage.NewMemoryStorage()
			taskService := service.NewTaskService(memoryStorage)

			_, err := taskService.CreateTask(
				tt.title,
				"Testing table-driven tests",
			)

			if tt.wantErr {
				if !errors.Is(err, model.ErrInvalidTask) {
					t.Fatalf("expected ErrInvalidTask, got %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}

		})
	}
}

func TestUpdateTask_TitleTooLong(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	taskService := service.NewTaskService(memoryStorage)

	task, err := taskService.CreateTask("Learn Go", "Study backend")
	if err != nil {
		t.Fatal(err)
	}

	longTitle := strings.Repeat("a", 101)

	err = taskService.UpdateTask(task.ID, &longTitle, nil)

	if !errors.Is(err, model.ErrInvalidTask) {
		t.Errorf("expecteed ErrInvalidTask, got %v", err)
	}
}

func TestUpdateTask_DescriptionTooLong(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	taskService := service.NewTaskService(memoryStorage)

	task, err := taskService.CreateTask("Learn Go", "Study backend")
	if err != nil {
		t.Fatal(err)
	}

	longDescription := strings.Repeat("a", 501)

	err = taskService.UpdateTask(task.ID, nil, &longDescription)

	if !errors.Is(err, model.ErrInvalidTask) {
		t.Errorf("expected ErrInvalidTask, got %v", err)
	}
}
