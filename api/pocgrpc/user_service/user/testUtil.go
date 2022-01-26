package user

import (
	"context"
	"errors"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/entities"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/pb"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (mock *mockService) CreateUser(ctx context.Context, user entities.User) (int64, error) {
	args := mock.Called(ctx, user)

	return int64(args.Int(0)), args.Error(1)
}

func (mock *mockService) GetUser(ctx context.Context, id int64) (entities.User, error) {
	args := mock.Called(ctx, id)

	return args.Get(0).(entities.User), args.Error(1)
}

type connectionMock struct {
	mock.Mock
}

func (conn *connectionMock) InsertUser(user entities.User) (int64, error) {
	args := conn.Called(user.Name, user.PwdHash, user.Age, user.ParentID)
	return int64(args.Int(0)), args.Error(1)
}

func (conn *connectionMock) GetUser(id int64) (entities.User, error) {
	args := conn.Called(id)

	return args.Get(0).(entities.User), args.Error(1)
}

func (conn connectionMock) GetUserByName(name string) (entities.User, error) {
	args := conn.Called(name)
	return args.Get(0).(entities.User), args.Error(1)
}

type mockErrorHandler struct {
	mock.Mock
}

func (mock *mockErrorHandler) TogRPCStatus(err error) *pb.Err {
	args := mock.Called(err)

	return args.Get(0).(*pb.Err)
}

func getGenericRepositoryError() error {
	return errors.New("generic repository error")
}

func getNewUser(options ...entities.InitializationOption) entities.User {
	user, _ := entities.NewUser("Test user", "Test password", 10, 0, options...)

	return user
}

func getNewgRPCUser() *pb.User {
	user := getNewUser()

	return &pb.User{
		Name:     user.Name,
		PwdHash:  user.PwdHash,
		Age:      uint32(user.Age),
		ParentId: user.ParentID,
	}
}
