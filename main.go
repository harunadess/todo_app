package main

import (
	"fmt"
	"html/template"
	"net/http"

	database "github.com/harunadess/todo_app/db"
	entities "github.com/harunadess/todo_app/entities"
	logger "github.com/harunadess/todo_app/logger"
)

var db database.DB

type ViewData struct {
	Lists        []entities.List
	Todos        []entities.Todo
	SelectedList int64
}

var templates map[string]*template.Template

func init() {
	if templates == nil {
		templates = GetTemplates()
	}

	db.Conn = database.OpenDbConnection()
	db.SetUp()
}

func registerStaticHandler() {
	http.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
}

func registerDefaultHandler() {
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("hit / endpoint")

		requestedPath := r.URL.Path
		logger.Info("requested path: ", requestedPath)
		if requestedPath != "/" {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		lists, err := db.GetAllLists()
		if err != nil {
			http.Error(w, "failed to get all lists", http.StatusInternalServerError)
		}

		if len(lists) == 0 {
			logger.Info("STUB: we don't have any lists, so we pretending for now.")

			viewData := ViewData{Lists: lists, Todos: make([]entities.Todo, 0), SelectedList: -1}
			tmpl := templates["index.html"]
			tmpl.Execute(w, viewData)
			return
		}

		todos, err := db.GetAllTodosInList(lists[0].ID)
		if err != nil {
			errStr := fmt.Sprintf("failed to get todos in list: %d", lists[0].ID)
			http.Error(w, errStr, http.StatusInternalServerError)
		}

		viewData := ViewData{Lists: lists, Todos: todos, SelectedList: lists[0].ID}
		tmpl := templates["index.html"]
		tmpl.Execute(w, viewData)
	})
}

func main() {
	RegisterTodoHandlers()
	RegisterListHandlers()
	registerStaticHandler()
	registerDefaultHandler()

	logger.Info("server is listening on :3000")
	logger.Fatal(http.ListenAndServe(":3000", nil))
}
