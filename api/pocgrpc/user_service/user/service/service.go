package service

import (
	"context"

	"github.com/go-kit/log"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/user/repository"
)

type Service interface {
	CreateUser(ctx context.Context, request CreateUserRequest) (CreateUserResponse, error)
}

type service struct {
	repository repository.Repository
	log        log.Logger
}

type CreateUserRequest struct {
	Name   string
	Pwd    string
	Age    uint8
	Parent int64
}
type CreateUserResponse struct {
	Id int
}

func NewService(repo repository.Repository, logger log.Logger) *service {
	return &service{
		repository: repo,
		log:        logger,
	}
}

//Creates an user if the information provided is valid
func (service *service) CreateUser(ctx context.Context, request CreateUserRequest) (*CreateUserResponse, error) {
	id, err := service.repository.InsertUser(repository.NewUser(request.Name, request.Pwd, request.Age, request.Parent))

	return &CreateUserResponse{
		Id: id,
	}, err
}
