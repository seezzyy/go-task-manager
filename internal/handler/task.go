package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"task-manager/internal/model"
	"task-manager/internal/service"
)

type TaskHandler struct {
	Service *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		Service: taskService,
	}
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Service.GetTasks()
	if err != nil {
		writeJSONEError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) GetTask(
	w http.ResponseWriter,
	r *http.Request,
) {
	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		writeJSONEError(
			w,
			http.StatusBadRequest,
			"Invalid task ID",
		)
		return
	}
	task, err := h.Service.GetTask(id)
	if errors.Is(err, model.ErrTaskNotFound) {
		writeJSONEError(
			w,
			http.StatusNotFound,
			"Task not found",
		)
		return
	}

	if err != nil {
		writeJSONEError(
			w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	Title       json.RawMessage `json:"title"`
	Description json.RawMessage `json:"description"`
}

func (h *TaskHandler) CreateTask(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request CreateTaskRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeJSONEError(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}
	task, err := h.Service.CreateTask(
		request.Title,
		request.Description,
	)

	if errors.Is(err, model.ErrInvalidTask) {
		writeValidationError(w, err)
		return
	}

	if err != nil {
		writeJSONEError(
			w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(task)

}

func (h *TaskHandler) CompleteTask(
	w http.ResponseWriter,
	r *http.Request,
) {
	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)

	if err != nil {
		writeJSONEError(
			w,
			http.StatusBadRequest,
			"Invalid task ID",
		)
		return
	}
	err = h.Service.CompleteTask(id)
	if errors.Is(err, model.ErrTaskNotFound) {
		writeJSONEError(
			w,
			http.StatusNotFound,
			"Task not found",
		)
		return
	}
	if err != nil {
		writeJSONEError(
			w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func (h *TaskHandler) DeleteTask(
	w http.ResponseWriter,
	r *http.Request,
) {
	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)

	if err != nil {
		writeJSONEError(
			w,
			http.StatusBadRequest,
			"Invalid task ID",
		)
		return
	}
	err = h.Service.DeleteTask(id)
	if errors.Is(err, model.ErrTaskNotFound) {
		writeJSONEError(
			w,
			http.StatusNotFound,
			"Task not found",
		)
		return
	}
	if err != nil {
		writeJSONEError(
			w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func (h *TaskHandler) UpdateTask(
	w http.ResponseWriter,
	r *http.Request,
) {
	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		writeJSONEError(
			w,
			http.StatusBadRequest,
			"Invalid task ID",
		)
		return
	}

	var request UpdateTaskRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeJSONEError(
			w,
			http.StatusBadRequest,
			"Invalid task",
		)
		return
	}

	var title *string

	if request.Title != nil {
		if string(request.Title) == "null" {
			writeJSONEError(
				w,
				http.StatusBadRequest,
				"Title cannot be null",
			)
			return
		}

		var value string

		err = json.Unmarshal(request.Title, &value)
		if err != nil {
			writeJSONEError(
				w,
				http.StatusBadRequest,
				"Invalid title",
			)
			return
		}

		title = &value
	}

	var description *string

	if request.Description != nil {
		if string(request.Description) == "null" {
			writeJSONEError(
				w,
				http.StatusBadRequest,
				"Description cannot be null",
			)
			return
		}

		var value string

		err = json.Unmarshal(request.Description, &value)
		if err != nil {
			writeJSONEError(
				w,
				http.StatusBadRequest,
				"Invalid description",
			)
			return
		}
		description = &value
	}

	err = h.Service.UpdateTask(
		id,
		title,
		description,
	)

	if errors.Is(err, model.ErrTaskNotFound) {
		writeJSONEError(
			w,
			http.StatusNotFound,
			"Task not found",
		)
		return
	}

	if errors.Is(err, model.ErrInvalidTask) {
		writeValidationError(w, err)
		return
	}

	if err != nil {
		writeJSONEError(
			w,
			http.StatusInternalServerError,
			"Invalid Server Error",
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSONEError(
	w http.ResponseWriter,
	status int,
	message string,
) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func writeValidationError(w http.ResponseWriter, err error) {
	message := "Invalid task"

	var validationErr *model.ValidationError

	if errors.As(err, &validationErr) {
		message = validationErr.Message
	}

	writeJSONEError(w, http.StatusBadRequest, message)
}
