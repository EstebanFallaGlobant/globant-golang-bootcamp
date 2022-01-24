package user

import (
	"context"

	"github.com/go-kit/log"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/entities"
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/error"
)

type Repository interface {
	InsertUser(user entities.User) (int64, error)
	GetUser(id int64) (entities.User, error)
	GetUserByName(name string) (entities.User, error)
}
type userInfoService struct {
	repository Repository
	log        log.Logger
}

func NewService(repo Repository, logger log.Logger) *userInfoService {
	return &userInfoService{
		repository: repo,
		log:        logger,
	}
}

// Creates an user if the information provided is valid.
//
// The users name and password must not be empty strings.
func (svc *userInfoService) CreateUser(ctx context.Context, user entities.User) (int64, error) {

	if userFound, _ := svc.repository.GetUserByName(user.Name); userFound.ID != 0 {
		return 0, svcerr.NewUserAlreadyExistError(userFound.Name, userFound.ID)
	}

	return svc.repository.InsertUser(user)
}

func (svc *userInfoService) GetUser(ctx context.Context, id int64) (entities.User, error) {
	if id < 1 {
		return entities.User{}, svcerr.NewInvalidArgumentsError(paramIDStr, ruleMsgID)
	}

	return svc.repository.GetUser(id)
}
