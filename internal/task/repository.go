package task

type Repository interface {
	GetAll() ([]Task, error)
	SaveAll([]Task) error
}
