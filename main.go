package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"slices"

	logger "github.com/harunadess/todo_app/logger"
)

// var db DB

var todoCount = 2
var todos = []Todo{
	{1, "Get eggs", false, true, 5},
	{2, "Get milk", true, false, 0},
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

func registerHandlers() {
	http.HandleFunc("POST /todos", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("hit /todo endpoint")

		name := r.PostFormValue("name")
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
		todo := Todo{id, name, false, hasCount, intCount}
		todos = append(todos, todo)

		logger.Info("added todo: ", todo)
		w.WriteHeader(http.StatusCreated)

		tmpl := templates["row.html"]
		tmpl.ExecuteTemplate(w, "row", todo)
	})

	http.HandleFunc("PUT /todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		intId, err := strconv.Atoi(id)
		if err != nil {
			logger.Error("failed to convert id to int: ", err)
			http.Error(w, "invalid id supplied: ", http.StatusBadRequest)
			return
		}

		name := r.PostFormValue("name")
		count := r.PostFormValue("count")
		done := r.PostFormValue("done") == "true"

		// get todo with id
		todoIdx := -1
		for i := range todos {
			if todos[i].ID == intId {
				todoIdx = i
				break
			}
		}
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
		todos[todoIdx].Count = intCount
		todos[todoIdx].Done = done
		logger.Info("updated item: ", todos[todoIdx])

		todo := todos[todoIdx]
		tmpl := templates["row.html"]
		tmpl.ExecuteTemplate(w, "row", todo)
	})

	http.HandleFunc("DELETE /todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		intId, err := strconv.Atoi(id)
		if err != nil {
			logger.Error("failed to convert id to int: ", err)
			http.Error(w, "invalid id supplied: ", http.StatusBadRequest)
			return
		}

		logger.Info("received DELETE for todo: ", intId)
		// get todo with id
		todoIdx := -1
		for i := range todos {
			if todos[i].ID == intId {
				todoIdx = i
				break
			}
		}
		if todoIdx == -1 {
			http.Error(w, "todo id not found", http.StatusNotFound)
			return
		}
		// remove todo from list
		todos = slices.Delete(todos, todoIdx, todoIdx+1)
	})

	http.HandleFunc("GET /todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		intId, err := strconv.Atoi(id)
		if err != nil {
			logger.Error("failed to convert id to int: ", err)
			http.Error(w, "invalid id supplied: ", http.StatusBadRequest)
			return
		}

		// get todo with id
		todoIdx := -1
		for i := range todos {
			if todos[i].ID == intId {
				todoIdx = i
				break
			}
		}
		if todoIdx == -1 {
			http.Error(w, "todo id not found", http.StatusNotFound)
			return
		}

		todo := todos[todoIdx]
		tmpl := templates["row.html"]
		tmpl.ExecuteTemplate(w, "row", todo)
	})

	http.HandleFunc("GET /todos/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		intId, err := strconv.Atoi(id)
		if err != nil {
			logger.Error("failed to convert id to int: ", err)
			http.Error(w, "invalid id supplied: ", http.StatusBadRequest)
			return
		}

		// get todo with id
		todoIdx := -1
		for i := range todos {
			if todos[i].ID == intId {
				todoIdx = i
				break
			}
		}
		if todoIdx == -1 {
			http.Error(w, "todo id not found", http.StatusNotFound)
			return
		}

		todo := todos[todoIdx]
		tmpl := templates["edit-item.html"]
		tmpl.ExecuteTemplate(w, "edit-item", todo)
	})

	http.HandleFunc("GET /todos", func(w http.ResponseWriter, r *http.Request) {
		content, err := json.Marshal(todos)
		if err != nil {
			logger.Fatal("unable to encode json: ", err)
		}

		w.Header().Add("content-type", "application/json")
		w.Write(content)
	})
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
	// tmpl := template.Must(CompileTemplates())
	registerHandlers()
	registerStaticHandler()
	registerDefaultHandler()

	logger.Info("server is listening on :3000")
	http.ListenAndServe(":3000", nil)
}
