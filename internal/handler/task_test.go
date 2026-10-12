package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"task-manager/internal/handler"
	"task-manager/internal/model"
	"task-manager/internal/service"
	"task-manager/internal/storage"
)

func TestGetTask_NotFound(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	taskService := service.NewTaskService(memoryStorage)
	taskHandler := handler.NewTaskHandler(taskService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks/999",
		nil,
	)

	req.SetPathValue("id", "999")

	rec := httptest.NewRecorder()

	taskHandler.GetTask(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d,got %d", http.StatusNotFound, rec.Code)
	}

	body := rec.Body.String()

	expectedBody := `{"error":"Task not found"}` + "\n"

	if body != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, body)
	}
}

func TestGetTask_Success(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	taskService := service.NewTaskService(memoryStorage)
	taskHandler := handler.NewTaskHandler(taskService)

	createdTask, err := taskService.CreateTask(
		"Learn Go",
		"Study HTTP testing",
	)

	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks/1",
		nil,
	)

	req.SetPathValue("id", strconv.Itoa(createdTask.ID))

	rec := httptest.NewRecorder()

	taskHandler.GetTask(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var receivedTask model.Task

	err = json.NewDecoder(rec.Body).Decode(&receivedTask)

	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if receivedTask.ID != createdTask.ID {
		t.Errorf("expected %d, got %d", createdTask.ID, receivedTask.ID)
	}

	if receivedTask.Title != "Learn Go" {
		t.Errorf("expected title %q, got %q", "Learn Go", receivedTask.Title)
	}

	if receivedTask.Description != "Study HTTP testing" {
		t.Errorf("expected description %q, got %q", "Study HTTP testing", receivedTask.Description)
	}

	if receivedTask.Status != "new" {
		t.Errorf("expected status %q, got %q", "new", receivedTask.Status)
	}

}

func TestCreateTask_Success(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	taskService := service.NewTaskService(memoryStorage)
	taskHandler := handler.NewTaskHandler(taskService)

	body := `{"title":"Learn Docker","description":"Study containers"}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(body),
	)
	req.Header.Set("Content-type", "application/json")

	rec := httptest.NewRecorder()

	taskHandler.CreateTask(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var createdTask model.Task

	err := json.NewDecoder(rec.Body).Decode(&createdTask)

	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if createdTask.ID != 1 {
		t.Errorf("expected ID %d, got %d", 1, createdTask.ID)
	}

	if createdTask.Title != "Learn Docker" {
		t.Errorf("expected title %q, got %q", "Learn Docker", createdTask.Title)
	}

	if createdTask.Description != "Study containers" {
		t.Errorf("expected description %q, got %q", "Study containers", createdTask.Description)
	}

	if createdTask.Status != "new" {
		t.Errorf("expected status %q, got %q", "new", createdTask.Status)
	}

	savedTask, err := taskService.GetTask(createdTask.ID)

	if err != nil {
		t.Fatalf("failed to get save task: %v", err)
	}

	if savedTask.ID != createdTask.ID {
		t.Errorf("expected saved ID %d, got %d", createdTask.ID, savedTask.ID)
	}
}

func TestCreatedTask_EmptyTitle(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	taskService := service.NewTaskService(memoryStorage)
	taskHandler := handler.NewTaskHandler(taskService)

	body := `{"title":"","description":"Some description"}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(body),
	)

	req.Header.Set("Content-type", "application/json")

	rec := httptest.NewRecorder()

	taskHandler.CreateTask(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	responseBody := rec.Body.String()

	expectedBody := "{\"error\":\"название пустое\"}\n"

	if responseBody != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, responseBody)
	}
}
