package storage

import (
	"database/sql"
	"errors"
	"task-manager/internal/model"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStorage struct {
	DB *sql.DB
}

func NewPostgresStorage(connectionString string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	postgresStorage := &PostgresStorage{
		DB: db,
	}
	return postgresStorage, nil
}

func (ps *PostgresStorage) Add(task model.Task) (int, error) {
	var id int

	query := `
	INSERT INTO tasks (title, description, status)
    VALUES ($1, $2, $3)
	RETURNING id
	`
	row := ps.DB.QueryRow(
		query,
		task.Title,
		task.Description,
		task.Status,
	)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (ps *PostgresStorage) GetByID(id int) (model.Task, error) {
	var task model.Task
	query := `
	SELECT id, title, description, status
	FROM tasks
	WHERE id = $1
	`
	row := ps.DB.QueryRow(query, id)
	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, model.ErrTaskNotFound
	}
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (ps *PostgresStorage) Update(task model.Task) error {
	query := `
	UPDATE tasks
	SET title = $1,
    	description = $2,
    	status = $3
	WHERE id = $4
	`

	result, err := ps.DB.Exec(
		query,
		task.Title,
		task.Description,
		task.Status,
		task.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return model.ErrTaskNotFound
	}

	return nil
}

func (ps *PostgresStorage) Delete(id int) error {
	query := `
	DELETE FROM tasks
	WHERE id = $1
	`
	result, err := ps.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrTaskNotFound
	}
	return nil
}

func (ps *PostgresStorage) List() ([]model.Task, error) {
	query := `
	SELECT id, title, description, status
	FROM tasks
	ORDER BY id
	`
	rows, err := ps.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := []model.Task{}
	for rows.Next() {
		var task model.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
