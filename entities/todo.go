package entities

type Todo struct {
	ID          int64
	ListID      int64
	Name        string
	Description string
	Done        bool
	HasCount    bool
	Count       int
}
