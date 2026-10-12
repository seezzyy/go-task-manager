//go:build integration

package storage_test

import (
	"errors"
	"testing"

	"task-manager/internal/model"
	"task-manager/internal/storage"
)

func TestPostgresStorage_Add(t *testing.T) {
	connectionString := "postgres://postgres:postgres@localhost:5434/taskmanager_test?sslmode=disable"

	ps, err := storage.NewPostgresStorage(connectionString)
	if err != nil {
		t.Fatalf("Не удалось подключиться к PostgresSQL: %v", err)
	}

	t.Cleanup(func() {
		ps.DB.Close()
	})

	task := model.Task{
		Title:       "Integration Test",
		Description: "Testing PostgresSQL",
		Status:      "new",
	}

	id, err := ps.Add(task)

	if err != nil {
		t.Fatalf("Ошибка при добавлении задачи: %v", err)
	}

	t.Cleanup(func() {
		if err := ps.Delete(id); err != nil {
			t.Errorf("Не удалось удалить тестовую задачу: %v", err)
		}
	})

	if id <= 0 {
		t.Fatalf("Ожидали положительный ID, получили %d", id)
	}

	savedTask, err := ps.GetByID(id)

	if err != nil {
		t.Fatalf("Не удалось получить добавленную задачу: %v", err)
	}

	if savedTask.Title != task.Title {
		t.Errorf(
			"Ожидали название %q, получили %q",
			task.Title,
			savedTask.Title,
		)
	}

	if savedTask.Description != task.Description {
		t.Errorf(
			"Ожидали описание %q, получили %q",
			task.Description,
			savedTask.Description,
		)
	}

	if savedTask.Status != task.Status {
		t.Errorf(
			"Ожидали статус %q, получили %q",
			task.Status,
			savedTask.Status,
		)
	}
}

func TestPostgresStorage_GetByID(t *testing.T) {
	connectionString := "postgres://postgres:postgres@localhost:5434/taskmanager_test?sslmode=disable"

	ps, err := storage.NewPostgresStorage(connectionString)
	if err != nil {
		t.Fatalf("Ошибка подключения: %v", err)
	}

	t.Cleanup(func() {
		ps.DB.Close()
	})

	task := model.Task{
		Title:       "Test GetByID",
		Description: "Checking GetByID",
		Status:      "new",
	}

	id, err := ps.Add(task)

	if err != nil {
		t.Fatalf("Ошибка при добавлении задачи: %v", err)
	}

	t.Cleanup(func() {
		if err := ps.Delete(id); err != nil {
			t.Errorf("Не удалось удалить тестовую задачу: %v", err)
		}
	})

	if id <= 0 {
		t.Fatalf("Ожидали положительное число: %d", id)
	}

	savedTask, err := ps.GetByID(id)

	if err != nil {
		t.Fatalf("Не ыдалось получить задачу: %v", err)
	}

	if savedTask.Title != task.Title {
		t.Errorf("Название задачи не совпадает")
	}

	if savedTask.ID != id {
		t.Errorf(
			"Ожидали ID %d, получили %d",
			id,
			savedTask.ID,
		)
	}

	if savedTask.Description != task.Description {
		t.Errorf(
			"Ожидали описание %q, получили %q",
			task.Description,
			savedTask.Description,
		)
	}

	if savedTask.Status != task.Status {
		t.Errorf(
			"Ожидали статус %q, получили %q",
			task.Status,
			savedTask.Status,
		)
	}

	_, err = ps.GetByID(-1)

	if !errors.Is(err, model.ErrTaskNotFound) {
		t.Errorf("Ожидали ErrTaskNotFound, получили %v", err)
	}
}

func TestPostgresStorage_Update(t *testing.T) {
	connectionString := "postgres://postgres:postgres@localhost:5434/taskmanager_test?sslmode=disable"

	ps, err := storage.NewPostgresStorage(connectionString)
	if err != nil {
		t.Fatalf("Ошибка подключения: %v", err)
	}

	t.Cleanup(func() {
		ps.DB.Close()
	})

	task := model.Task{
		Title:       "Original Title",
		Description: "Original Description",
		Status:      "new",
	}

	id, err := ps.Add(task)
	if err != nil {
		t.Fatalf("Ошибка добавления: %v", err)
	}

	t.Cleanup(func() {
		if err := ps.Delete(id); err != nil {
			t.Errorf("Ошибка очистки: %v", err)
		}
	})

	task.ID = id
	task.Title = "Updated Title"
	task.Description = "Updated Description"
	task.Status = "completed"

	err = ps.Update(task)
	if err != nil {
		t.Fatalf("Ошибка обновления: %v", err)
	}

	updatedTask, err := ps.GetByID(id)
	if err != nil {
		t.Fatalf("Ошибка получения: %v", err)
	}

	if updatedTask != task {
		t.Errorf("Ожидали %+v, получили %+v", task, updatedTask)
	}

	missingTask := model.Task{
		ID:     -1,
		Title:  "Missing Task",
		Status: "new",
	}

	err = ps.Update(missingTask)

	if !errors.Is(err, model.ErrTaskNotFound) {
		t.Errorf("Ожидали ErrTaskNotFound, получили %v", err)
	}
}

func TestPostgresStorage_Delete(t *testing.T) {
	connectionString := "postgres://postgres:postgres@localhost:5434/taskmanager_test?sslmode=disable"

	ps, err := storage.NewPostgresStorage(connectionString)
	if err != nil {
		t.Fatalf("Ошибка подключения: %v", err)
	}

	t.Cleanup(func() {
		ps.DB.Close()
	})

	task := model.Task{
		Title:       "Delete Test",
		Description: "Testing Delete",
		Status:      "new",
	}

	id, err := ps.Add(task)
	if err != nil {
		t.Fatalf("Ошибка добавления: %v", err)
	}

	// Подстраховка: если тест прервётся до удаления,
	// попробуем очистить запись после завершения.
	t.Cleanup(func() {
		_, err := ps.DB.Exec("DELETE FROM tasks WHERE id = $1", id)
		if err != nil {
			t.Errorf("Ошибка очистки: %v", err)
		}
	})

	err = ps.Delete(id)
	if err != nil {
		t.Fatalf("Ошибка удаления: %v", err)
	}

	_, err = ps.GetByID(id)

	if !errors.Is(err, model.ErrTaskNotFound) {
		t.Errorf("Ожидали ErrTaskNotFound, получили %v", err)
	}

	err = ps.Delete(-1)

	if !errors.Is(err, model.ErrTaskNotFound) {
		t.Errorf("Ожидали ErrTaskNotFound, получили %v", err)
	}
}

func TestPostgresStorage_List(t *testing.T) {
	connectionString := "postgres://postgres:postgres@localhost:5434/taskmanager_test?sslmode=disable"

	ps, err := storage.NewPostgresStorage(connectionString)
	if err != nil {
		t.Fatalf("Ошибка подключения: %v", err)
	}

	t.Cleanup(func() {
		ps.DB.Close()
	})

	task1 := model.Task{
		Title:  "List Task 1",
		Status: "new",
	}

	task2 := model.Task{
		Title:  "List Task 2",
		Status: "completed",
	}

	id1, err := ps.Add(task1)
	if err != nil {
		t.Fatalf("Ошибка добавления первой задачи: %v", err)
	}

	t.Cleanup(func() {
		if err := ps.Delete(id1); err != nil {
			t.Errorf("Ошибка очистки первой задачи: %v", err)
		}
	})

	id2, err := ps.Add(task2)
	if err != nil {
		t.Fatalf("Ошибка добавления второй задачи: %v", err)
	}

	t.Cleanup(func() {
		if err := ps.Delete(id2); err != nil {
			t.Errorf("Ошибка очистки второй задачи: %v", err)
		}
	})

	tasks, err := ps.List()
	if err != nil {
		t.Fatalf("Ошибка получения списка: %v", err)
	}

	found1 := false
	found2 := false

	for _, task := range tasks {
		if task.ID == id1 && task.Title == task1.Title && task.Status == task1.Status {
			found1 = true
		}

		if task.ID == id2 && task.Title == task2.Title && task.Status == task2.Status {
			found2 = true
		}
	}

	if !found1 || !found2 {
		t.Errorf("Не все добавленные задачи найдены: found1=%v, found2=%v",
			found1, found2)
	}
}
