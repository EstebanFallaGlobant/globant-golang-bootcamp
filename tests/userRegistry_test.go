package tests

import (
	"testing"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/Part3/userregistry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRegister struct {
	mock.Mock
}

func (register *mockRegister) Register(email string, pass string) (int, error) {
	args := register.Called(email, pass)

	return args.Int(0), args.Error(1)
}

func Test_userRegister_ValidEmailAndPass(t *testing.T) {
	email, pass := "test@mail.com", "password"

	register := new(mockRegister)
	register.On("Register", email, pass).Return(200, nil)

	userRegister := userregistry.UserRegister{Reg: register}

	userId, err := userRegister.Register(email, pass)

	assert.Nil(t, err)
	assert.EqualValues(t, 200, userId)
}

func Test_userRegister_InvalidEmail(t *testing.T) {
	email, pass := "testmail.com", "password"

	register := new(mockRegister)
	register.On("Register", email, pass).Return(200, nil)

	userRegister := userregistry.UserRegister{Reg: register}

	userId, err := userRegister.Register(email, pass)

	assert.EqualError(t, err, "invalid e-mail")
	assert.NotEqualValues(t, 200, userId)
}

func Test_userRegister_InvalidPassword(t *testing.T) {
	email, pass := "test@mail.com", "pass"

	register := new(mockRegister)
	register.On("Register", email, pass).Return(200, nil)

	userRegister := userregistry.UserRegister{Reg: register}

	userId, err := userRegister.Register(email, pass)

	assert.EqualError(t, err, "invalid password")
	assert.NotEqualValues(t, 200, userId)
}
