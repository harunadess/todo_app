package db

import (
	entities "github.com/harunadess/todo_app/entities"
	logger "github.com/harunadess/todo_app/logger"
)

func (db DB) SetUpTodosTable() error {
	_, err := db.Conn.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY,
			list_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			done INTEGER NOT NULL DEFAULT 0,
			has_count INTEGER NOT NULL DEFAULT 0,
			count INTEGER NOT NULL DEFAULT 0
		);
	`)

	if err != nil {
		logger.Fatal("error creating todos table: ", err)
		return err
	}

	return nil
}

func (db DB) CreateTodo(todo entities.Todo) (int64, error) {
	sql := `INSERT INTO todos (list_id, name, description, done, has_count, count)
			VALUES (?, ?, ?, ?, ?, ?);`
	result, err := db.Conn.Exec(sql, todo.ListID, todo.Name, todo.Description, todo.Done, todo.HasCount, todo.Count)
	if err != nil {
		logger.Error("failed to create todo: ", err)
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("failed to get id of inserted todo: ", err)
		return -1, nil
	}

	logger.Info("created todo: ", id)
	return id, nil
}

func (db DB) UpdateTodo(id int64, name string, description string, hasCount bool, count int) error {
	sql := `UPDATE todos
			SET name = ?, description = ?, has_count = ?, count = ?
			WHERE id = ?;`
	result, err := db.Conn.Exec(sql, name, description, hasCount, count, id)
	if err != nil {
		logger.Error("failed to update todo: ", err)
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		logger.Error("failed to get number of updated rows: ", err)
		return nil
	}

	logger.Info("updated todo with id: ", id, " rows affected: ", affected)
	return nil
}

func (db DB) ToggleTodoDone(id int64, done bool) error {
	sql := `UPDATE todos
			SET done = ?
			WHERE id = ?;`
	result, err := db.Conn.Exec(sql, done, id)
	if err != nil {
		logger.Error("failed to set done on todo: ", err)
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		logger.Error("failed to get number of updated rows: ", err)
		return nil
	}

	logger.Info("set done on todo with id: ", id, " rows affected: ", affected)
	return nil
}

func (db DB) DeleteTodo(id int64) error {
	sql := "DELETE FROM todos WHERE id = ?;"
	result, err := db.Conn.Exec(sql, id)
	if err != nil {
		logger.Error("failed to delete todo: ", err)
		return nil
	}

	affected, err := result.RowsAffected()
	if err != nil {
		logger.Error("failed to get number of updated rows: ", err)
		return nil
	}

	logger.Info("deleted todo with id: ", id, " rows affected: ", affected)
	return nil
}

func (db DB) GetAllTodosInList(listId int64) ([]entities.Todo, error) {
	sql := "SELECT * FROM todos WHERE list_id = ? ORDER BY id;"
	rows, err := db.Conn.Query(sql, listId)
	if err != nil {
		logger.Error("failed to get all todos: ", err)
		return nil, err
	}
	defer rows.Close()

	var todos []entities.Todo
	for rows.Next() {
		todo := &entities.Todo{}
		err := rows.Scan(&todo.ID, &todo.ListID, &todo.Name, &todo.Description, &todo.Done, &todo.HasCount, &todo.Count)
		if err != nil {
			logger.Error("failed to convert list: ", err)
			return nil, err
		}
		todos = append(todos, *todo)
	}

	return todos, nil
}

func (db DB) GetTodo(id int64) (*entities.Todo, error) {
	sql := "SELECT * FROM todos WHERE id = ?;"
	row := db.Conn.QueryRow(sql, id)

	todo := &entities.Todo{}
	err := row.Scan(&todo.ID, &todo.ListID, &todo.Name, &todo.Description, &todo.Done, &todo.HasCount, &todo.Count)
	if err != nil {
		logger.Error("failed to convert todo: ", err)
		return nil, err
	}

	return todo, nil
}
