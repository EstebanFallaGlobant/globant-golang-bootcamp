package user

import (
	"context"
	"errors"

	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (mock *mockService) CreateUser(ctx context.Context, user User) (int64, error) {
	args := mock.Called(ctx, user)

	return int64(args.Int(0)), args.Error(1)
}

func (mock *mockService) GetUser(ctx context.Context, id int64) (User, error) {
	args := mock.Called(ctx, id)

	return args.Get(0).(User), args.Error(1)
}

type connectionMock struct {
	mock.Mock
}

func (connection *connectionMock) InsertUser(user User) (int64, error) {
	args := connection.Called(user.Name, user.PwdHash, user.Age, user.Parent)
	return int64(args.Int(0)), args.Error(1)
}

func (connection *connectionMock) GetUser(id int64) (User, error) {
	args := connection.Called(id)

	return args.Get(0).(User), args.Error(1)
}

func getGenericRepositoryError() error {
	return errors.New("generic repository error")
}

func getNewUser() User {
	user, _ := NewUser("Test user", "Test password", 10, 0)

	return user
}
