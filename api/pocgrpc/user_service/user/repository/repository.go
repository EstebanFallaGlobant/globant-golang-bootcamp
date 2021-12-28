package repository

type Repository interface {
	InsertUser(user User) (int, error)
}
