package user

import (
	"context"
	"fmt"

	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (mock *mockService) CreateUser(ctx context.Context, user User) (int64, error) {
	args := mock.Called(ctx, user)

	fmt.Printf("Returns: %d, %v\n", args.Int(0), args.Error(1))

	return int64(args.Int(0)), args.Error(1)
}
