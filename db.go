package main

import (
	"database/sql"

	logger "github.com/harunadess/todo_app/logger"

	_ "github.com/glebarez/go-sqlite"
)

type DB struct {
	conn *sql.DB
}

func (db DB) Connect() *sql.DB {
	conn, err := sql.Open("sqlite", "./todo.db")
	if err != nil {
		logger.Fatal("error opening db: ", err)
	}

	logger.Info("successfully opened db connection")

	return conn
}

func (db DB) SetUp() {
	_, err := db.conn.Query(`
		CREATE TABLE IF NOT EXISTS todo_list (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			done INTEGER NOT NULL DEFAULT 0,
			has_count INTEGER NOT NULL DEFAULT 0,
			count INTEGER NOT NULL DEFAULT 0
		);
	`)

	if err != nil {
		logger.Fatal("error creating todo_list table: ", err)
	}
}

func (db DB) GetAll() ([]Todo, error) {
	rows, err := db.conn.Query("SELECT * FROM todo_list ORDER BY id;")
	if err != nil {
		return []Todo{}, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		todo := &Todo{}
		err := rows.Scan(&todo.ID, &todo.Name, &todo.Done, &todo.HasCount, &todo.Count)
		if err != nil {
			logger.Error("error converting row: ", err)
			continue
		}

		todos = append(todos, *todo)
	}

	return todos, nil
}
