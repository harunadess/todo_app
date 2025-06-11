package main

type Todo struct {
	ID          int
	Name        string
	Description string
	Done        bool
	HasCount    bool
	Count       int
}

type List struct {
	ID            int
	Name          string
	CompletedDate string
}

type ListWithTodos struct {
	List  List
	Todos []Todo
}
