package db

import (
	"database/sql"

	logger "github.com/harunadess/todo_app/logger"

	_ "github.com/glebarez/go-sqlite"
)

type DB struct {
	Conn *sql.DB
}

func OpenDbConnection() *sql.DB {
	conn, err := sql.Open("sqlite", "./todo.db")
	if err != nil {
		logger.Fatal("error opening db: ", err)
	}

	logger.Info("successfully opened db connection")

	return conn
}

func (db DB) SetUp() {
	err := db.SetUpTodosTable()
	if err != nil {
		logger.Fatal("failed db setup: ", err)
	}
	err = db.SetUpListsTable()
	if err != nil {
		logger.Fatal("failed db setup: ", err)
	}

	// temp
}
