package entities

type List struct {
	ID            int64
	Name          string
	Completed     bool
	CompletedDate string
}

type ListWithTodos struct {
	List  List
	Todos []Todo
}
