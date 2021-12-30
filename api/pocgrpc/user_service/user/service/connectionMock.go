package service

import (
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/user/repository"
	"github.com/stretchr/testify/mock"
)

type connectionMock struct {
	mock.Mock
}

func (connection *connectionMock) InsertUser(user repository.User) (int64, error) {
	args := connection.Called(user.Name, user.PwdHash, user.Age, user.Parent)
	return int64(args.Int(0)), args.Error(1)
}
