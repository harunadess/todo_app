package db

import (
	logger "github.com/harunadess/todo_app/logger"

	entities "github.com/harunadess/todo_app/entities"
)

// set up [x]
// get all lists [x]
// create list [x]
// update list [ ]
// delete list [ ]

func (db DB) SetUpListsTable() error {
	_, err := db.Conn.Exec(`
		CREATE TABLE IF NOT EXISTS lists (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			completed INTEGER NOT NULL DEFAULT 0,
			completed_date TEXT
		);
	`)

	if err != nil {
		logger.Fatal("error creating lists table: ", err)
		return err
	}

	return nil
}

func (db DB) GetAllLists() ([]entities.List, error) {
	sql := "SELECT * FROM lists ORDER BY id;"
	rows, err := db.Conn.Query(sql)
	if err != nil {
		logger.Error("failed to get all lists: ", err)
		return nil, err
	}
	defer rows.Close()

	var lists []entities.List
	for rows.Next() {
		list := &entities.List{}
		err := rows.Scan(&list.ID, &list.Name, &list.Completed, &list.CompletedDate)
		if err != nil {
			logger.Error("failed to convert list: ", err)
			return nil, err
		}
		lists = append(lists, *list)
	}

	return lists, nil
}

func (db DB) CreateList(name string) (int64, error) {
	sql := `INSERT INTO lists (name)
			VALUES (?);`
	result, err := db.Conn.Exec(sql, name)
	if err != nil {
		logger.Error("failed to create list: ", err)
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("failed to get id of inserted list: ", err)
		return -1, nil
	}

	logger.Info("created list: ", id)
	return id, nil
}
