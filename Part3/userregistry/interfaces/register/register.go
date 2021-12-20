package register

type Register interface {
	//Save an user to the persistance layer
	Register(email string, pass string) (int, error)
}
