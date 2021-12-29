package service

import (
	"context"

	"github.com/go-kit/log"

	apiErr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/user/errors"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/user/repository"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/util"
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
	var err error = nil
	var params []string
	var response CreateUserResponse

	if util.IsEmptyString(request.Name) {
		params = append(params, "name")
	}

	if util.IsEmptyString(request.Pwd) {
		params = append(params, "password")
	}

	if id, v := 0, len(params) > 0; v {
		err = apiErr.NewArgumentsRequiredError(params...)
	} else if age := request.Age; age < 1 || age > 150 {
		err = apiErr.NewInvalidArgumentsError("age", "between 1 and 150")
	} else {
		id, err = service.repository.InsertUser(
			repository.NewUser(
				request.Name,
				request.Pwd,
				request.Age,
				request.Parent))

		response.Id = id
	}

	return &response, err
}
