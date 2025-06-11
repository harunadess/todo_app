package main

import (
	"net/http"
	"slices"
	"strconv"

	logger "github.com/harunadess/todo_app/logger"
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

	todoCount += 1
	id := todoCount
	todo := Todo{id, name, desc, false, hasCount, intCount}
	todos = append(todos, todo)

	logger.Info("added todo: ", todo)
	w.WriteHeader(http.StatusCreated)

	tmpl := templates["row.html"]
	tmpl.ExecuteTemplate(w, "row", todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		logger.Error("failed to convert id to int: ", err)
		http.Error(w, "invalid id supplied: ", http.StatusBadRequest)
		return
	}

	name := r.PostFormValue("name")
	desc := r.PostFormValue("description")
	count := r.PostFormValue("count")
	done := r.PostFormValue("done") == "true"

	todoIdx := getIdxOfTodo(intId)
	if todoIdx == -1 {
		http.Error(w, "todo id not found", http.StatusNotFound)
		return
	}

	intCount := 0
	if todos[todoIdx].HasCount {
		i, err := strconv.Atoi(count)
		if err != nil {
			logger.Error("failed to convert count to int: ", err)
			http.Error(w, "invalid count supplied", http.StatusBadRequest)
			return
		}

		intCount = i
	}

	todos[todoIdx].Name = name
	todos[todoIdx].Description = desc
	todos[todoIdx].Count = intCount
	todos[todoIdx].Done = done
	logger.Info("updated item: ", todos[todoIdx])

	todo := todos[todoIdx]
	tmpl := templates["row.html"]
	tmpl.ExecuteTemplate(w, "row", todo)
}

func toggleDone(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	done := r.URL.Query().Get("done") == "true"
	intId, err := strconv.Atoi(id)
	if err != nil {
		logger.Error("failed to convert id to int: ", err)
		http.Error(w, "invalid id supplied: ", http.StatusBadRequest)
		return
	}

	todoIdx := getIdxOfTodo(intId)
	if todoIdx == -1 {
		http.Error(w, "todo id not found", http.StatusNotFound)
		return
	}

	todos[todoIdx].Done = done
	logger.Info("updated item: ", todos[todoIdx])

	todo := todos[todoIdx]
	tmpl := templates["row.html"]
	tmpl.ExecuteTemplate(w, "row", todo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		logger.Error("failed to convert id to int: ", err)
		http.Error(w, "invalid id supplied: ", http.StatusBadRequest)
		return
	}

	logger.Info("received DELETE for todo: ", intId)

	todoIdx := getIdxOfTodo(intId)
	if todoIdx == -1 {
		http.Error(w, "todo id not found", http.StatusNotFound)
		return
	}
	// remove todo from list
	todos = slices.Delete(todos, todoIdx, todoIdx+1)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		logger.Error("failed to convert id to int: ", err)
		http.Error(w, "invalid id supplied: ", http.StatusBadRequest)
		return
	}

	todoIdx := getIdxOfTodo(intId)
	if todoIdx == -1 {
		http.Error(w, "todo id not found", http.StatusNotFound)
		return
	}

	todo := todos[todoIdx]
	tmpl := templates["row.html"]
	tmpl.ExecuteTemplate(w, "row", todo)
}

func getEditView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		logger.Error("failed to convert id to int: ", err)
		http.Error(w, "invalid id supplied: ", http.StatusBadRequest)
		return
	}

	todoIdx := getIdxOfTodo(intId)
	if todoIdx == -1 {
		http.Error(w, "todo id not found", http.StatusNotFound)
		return
	}

	todo := todos[todoIdx]
	tmpl := templates["edit-item.html"]
	tmpl.ExecuteTemplate(w, "edit-item", todo)
}

// temp function to get idx of todo from static list
func getIdxOfTodo(id int) int {
	todoIdx := -1
	for i := range todos {
		if todos[i].ID == id {
			todoIdx = i
			break
		}
	}

	return todoIdx
}
