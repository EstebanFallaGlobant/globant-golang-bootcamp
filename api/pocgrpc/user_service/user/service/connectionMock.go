package service

import (
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/user/repository"
	"github.com/stretchr/testify/mock"
)

type connectionMock struct {
	mock.Mock
}

func (connection *connectionMock) InsertUser(user repository.User) (int, error) {
	args := connection.Called(user.Name, user.PwdHash, user.Age, user.Parent)
	return args.Int(0), args.Error(1)
}
