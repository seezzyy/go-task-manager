package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"task-manager/internal/config"
	"task-manager/internal/handler"
	"task-manager/internal/middleware"
	"task-manager/internal/service"
	"task-manager/internal/storage"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Task Manager API is running")
}

func main() {
	logHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	cfg, err := config.Load()

	if err != nil {
		slog.Error("Ошибка загрузки конфигурации:", "error", err)
		return
	}

	connectionString := cfg.DatabaseURL

	postgresStorage, err := storage.NewPostgresStorage(connectionString)
	if err != nil {
		slog.Error("Ошибка подключения к PostgreSQL:", "error", err)
		return
	}

	defer postgresStorage.DB.Close()

	taskService := service.NewTaskService(postgresStorage)
	taskHandler := handler.NewTaskHandler(taskService)

	http.HandleFunc("GET /tasks", taskHandler.GetTasks)
	http.HandleFunc("GET /tasks/{id}", taskHandler.GetTask)
	http.HandleFunc("POST /tasks", taskHandler.CreateTask)
	http.HandleFunc("PATCH /tasks/{id}", taskHandler.UpdateTask)
	http.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)

	http.HandleFunc("/", homeHandler)

	server := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           middleware.Logging(http.DefaultServeMux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		slog.Info("HTTP-сервер запущен", "port", cfg.HTTPPort)

		err := server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case <-ctx.Done():
		slog.Info("Получен сигнал остановки")

	case err := <-serverErrors:
		slog.Error("Ошибка HTTP-сервера:", "error", err)
		return
	}

	stop()

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		cfg.ShutdownTimeout,
	)
	defer cancel()

	slog.Info("Останавливаем HTTP-сервер...")

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Graceful Shutdown завершился с ошибкой:", "error", err)

		if closeErr := server.Close(); closeErr != nil {
			slog.Error("Ошибка принудительного закрытия:", "error", closeErr)
		}
	}

	slog.Info("HTTP-сервер остановлен")
}
