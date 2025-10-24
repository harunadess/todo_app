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

func RegisterListHandlers() {
	http.HandleFunc("GET /lists/{id}/todos", getAllTodosInList)
	http.HandleFunc("POST /lists", addList)
	http.HandleFunc("PUT /lists/{id}", updateList)
	http.HandleFunc("DELETE /lists/{id}", deleteList)
	http.HandleFunc("GET /list/{id}/edit", getListEditView)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	logger.Info("hit /todo endpoint")

	name := r.PostFormValue("name")
	desc := r.PostFormValue("description")
	hasCount := r.PostFormValue("has-count") == "true"
	count := r.PostFormValue("count")

	listIdParam := r.URL.Query().Get("listId")
	if len(listIdParam) == 0 {
		logger.Error("listId was empty: ", listIdParam)
		http.Error(w, "listId was empty", http.StatusBadRequest)
		return
	}
	listId, err := strconv.ParseInt(listIdParam, 10, 64)
	if err != nil {
		logger.Error("failed to convert listId to int: ", err)
		http.Error(w, "invalid listId supplied", http.StatusBadRequest)
		return
	}

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

	todo := entities.Todo{ListID: listId, Name: name, Description: desc, Done: false, HasCount: hasCount, Count: intCount}
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

func getAllTodosInList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logger.Error("failed to convert id to int: ", err)
		http.Error(w, "invalid id supplied: ", http.StatusBadRequest)
		return
	}

	todos, err := db.GetAllTodosInList(intId)
	if err != nil {
		logger.Error("failed to get todo: ", err)
		http.Error(w, "failed to get todo", http.StatusInternalServerError)
		return
	}
	tmpl := templates["item-content.html"]
	viewData := ViewData{Lists: make([]entities.List, 0), Todos: todos, SelectedList: intId}
	tmpl.ExecuteTemplate(w, "item-content", viewData)
}

func addList(w http.ResponseWriter, r *http.Request) {
	logger.Info("hit /todo endpoint")

	name := r.PostFormValue("name")
	list := entities.List{Name: name, CompletedDate: ""}
	id, err := db.CreateList(list)
	if err != nil {
		logger.Error("failed to create list: ", err)
		http.Error(w, "failed to create list", http.StatusInternalServerError)
		return
	}
	list.ID = id

	logger.Info("added list: ", list)
	w.WriteHeader(http.StatusCreated)

	tmpl := templates["nav-row.html"]
	tmpl.ExecuteTemplate(w, "nav-row", list)
}

func updateList(w http.ResponseWriter, r *http.Request) {}

/*
On deleting of a list, need to fix this to return the entire view again (kind of)
which might result in more db hits so we know what lists we have

I don't really like that idea, so idk how to resolve this right now.
*/
func deleteList(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logger.Error("failed to convert id to int: ", err)
		http.Error(w, "invalid id supplied: ", http.StatusBadRequest)
		return
	}

	err = db.DeleteAllTodosInList(intId)
	if err != nil {
		logger.Error("failed to delete todos: ", err)
		http.Error(w, "failed to delete todos", http.StatusInternalServerError)
		return
	}

	err = db.DeleteList(intId)
	if err != nil {
		logger.Error("failed to delete list: ", err)
		http.Error(w, "failed to delete list", http.StatusInternalServerError)
		return
	}

	tmpl := templates["item-content.html"]
	viewData := ViewData{Lists: make([]entities.List, 0), Todos: make([]entities.Todo, 0), SelectedList: -1}
	tmpl.ExecuteTemplate(w, "item-content", viewData)
}

func getListEditView(w http.ResponseWriter, r *http.Request) {}
