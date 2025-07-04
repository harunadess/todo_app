package main

import (
	"net/http"
	"strconv"

	entities "github.com/harunadess/todo_app/entities"
	"github.com/harunadess/todo_app/logger"
)

func RegisterTodoHandlers() {
	http.HandleFunc("POST /todos", addTodo)
	http.HandleFunc("PUT /todos/{id}", updateTodo)
	http.HandleFunc("PUT /todos/{id}/toggledone", toggleDone)
	http.HandleFunc("DELETE /todos/{id}", deleteTodo)
	http.HandleFunc("GET /todos/{id}", getTodo)
	http.HandleFunc("GET /todos/{id}/edit", getEditView)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	logger.Info("hit /todo endpoint")

	name := r.PostFormValue("name")
	desc := r.PostFormValue("description")
	hasCount := r.PostFormValue("has-count") == "true"
	count := r.PostFormValue("count")

	intCount := 0
	if hasCount {
		i, err := strconv.Atoi(count)
		if err != nil {
			logger.Error("failed to convert count to int: ", err)
			http.Error(w, "invalid count supplied", http.StatusBadRequest)
			return
		}

		intCount = i
	}

	todo := entities.Todo{ListID: 1, Name: name, Description: desc, Done: false, HasCount: hasCount, Count: intCount}
	id, err := db.CreateTodo(todo)
	if err != nil {
		logger.Error("failed to create todo: ", err)
		http.Error(w, "failed to create todo", http.StatusInternalServerError)
		return
	}
	todo.ID = id

	logger.Info("added todo: ", todo)
	w.WriteHeader(http.StatusCreated)

	tmpl := templates["row.html"]
	tmpl.ExecuteTemplate(w, "row", todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logger.Error("failed to convert id to int: ", err)
		http.Error(w, "invalid id supplied: ", http.StatusBadRequest)
		return
	}

	name := r.PostFormValue("name")
	desc := r.PostFormValue("description")
	hasCount := r.PostFormValue("has-count") == "true"
	count := r.PostFormValue("count")

	intCount, err := strconv.Atoi(count)
	if err != nil {
		logger.Error("failed to convert count to int: ", err)
		http.Error(w, "invalid count supplied: ", http.StatusBadRequest)
		return
	}

	if !hasCount {
		intCount = 0
	}

	err = db.UpdateTodo(intId, name, desc, hasCount, intCount)
	if err != nil {
		logger.Error("failed to update todo: ", err)
		http.Error(w, "failed to update todo", http.StatusInternalServerError)
		return
	}

	todo, err := db.GetTodo(intId)
	if err != nil {
		logger.Error("failed to get todo: ", err)
		http.Error(w, "failed to get todo after update", http.StatusInternalServerError)
		return
	}

	tmpl := templates["row.html"]
	tmpl.ExecuteTemplate(w, "row", *todo)
}

func toggleDone(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	done := r.URL.Query().Get("done") == "true"
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logger.Error("failed to convert id to int: ", err)
		http.Error(w, "invalid id supplied: ", http.StatusBadRequest)
		return
	}

	err = db.ToggleTodoDone(intId, done)
	if err != nil {
		logger.Error("failed to toggle done: ", err)
		http.Error(w, "failed to toggle done for todo", http.StatusInternalServerError)
		return
	}
	logger.Info("updated todo: ", intId)

	todo, err := db.GetTodo(intId)
	if err != nil {
		logger.Error("failed to get todo: ", err)
		http.Error(w, "failed to get todo after toggle done", http.StatusInternalServerError)
		return
	}

	tmpl := templates["row.html"]
	tmpl.ExecuteTemplate(w, "row", *todo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logger.Error("failed to convert id to int: ", err)
		http.Error(w, "invalid id supplied: ", http.StatusBadRequest)
		return
	}

	err = db.DeleteTodo(intId)
	if err != nil {
		logger.Error("failed to delete todo: ", err)
		http.Error(w, "failed to delete todo", http.StatusInternalServerError)
		return
	}
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logger.Error("failed to convert id to int: ", err)
		http.Error(w, "invalid id supplied: ", http.StatusBadRequest)
		return
	}

	todo, err := db.GetTodo(intId)
	if err != nil {
		logger.Error("failed to get todo: ", err)
		http.Error(w, "failed to get todo", http.StatusInternalServerError)
		return
	}
	tmpl := templates["row.html"]
	tmpl.ExecuteTemplate(w, "row", *todo)
}

func getEditView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logger.Error("failed to convert id to int: ", err)
		http.Error(w, "invalid id supplied: ", http.StatusBadRequest)
		return
	}

	todo, err := db.GetTodo(intId)
	if err != nil {
		logger.Error("failed to get todo: ", err)
		http.Error(w, "failed to get todo", http.StatusInternalServerError)
		return
	}

	tmpl := templates["edit-item.html"]
	tmpl.ExecuteTemplate(w, "edit-item", *todo)
}
