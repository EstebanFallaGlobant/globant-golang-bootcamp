package userregistry

import (
	"errors"
	"regexp"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/Part3/userregistry/interfaces/register"
)

type UserRegister struct {
	Reg register.Register
}

// Checks the main and password and if they are valid saved the information. Returns the id of the saved user
func (userRegister *UserRegister) Register(email string, pass string) (int, error) {
	var status int
	var err error
	// Regular expression taken from: https://regex101.com/library/RzBwPX
	regxp := regexp.MustCompile(`^(?P<name>[a-zA-Z0-9.!#$%&'*+/=?^_ \x60{|}~-]+)@(?P<domain>[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)$`)

	if v1, v2 := regxp.MatchString(email), len(pass) >= 8; v1 && v2 {
		status, err = userRegister.Reg.Register(email, pass)
	} else if !v1 {
		status = 0
		err = errors.New("invalid e-mail")
	} else if !v2 {
		status = 0
		err = errors.New("invalid password")
	}

	return status, err
}
