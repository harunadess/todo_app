package main

import (
	"html/template"
	"net/http"

	logger "github.com/harunadess/todo_app/logger"
)

// var db DB

var todoCount = 3
var todos = []Todo{
	{1, "Get eggs", "Can get from store", false, true, 5},
	{2, "Get milk", "", true, false, 0},
	{3, "Get alexandrite", "Buy from Auriana, 50 poetics each", true, true, 50},
}

type ViewData struct {
	Todos []Todo
}

var templates map[string]*template.Template

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templates["index.html"] = template.Must(template.ParseFiles("templates/index.html", "templates/add-form.html", "templates/list.html", "templates/row.html"))
	templates["list.html"] = template.Must(template.ParseFiles("templates/list.html", "templates/row.html"))
	templates["row.html"] = template.Must(template.ParseFiles("templates/row.html"))
	templates["edit-item.html"] = template.Must(template.ParseFiles("templates/edit-item.html"))

	// db.conn = db.Connect()
	// db.SetUp()
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

		viewData := ViewData{Todos: todos}
		tmpl := templates["index.html"]
		tmpl.Execute(w, viewData)
	})
}

func main() {
	RegisterTodoHandlers()
	registerStaticHandler()
	registerDefaultHandler()

	logger.Info("server is listening on :3000")
	logger.Fatal(http.ListenAndServe(":3000", nil))
}
