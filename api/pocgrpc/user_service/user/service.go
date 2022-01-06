package user

import (
	"context"

	"github.com/go-kit/log"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/util"
)

type Repository interface {
	InsertUser(user User) (int64, error)
	GetUser(id int64) (User, error)
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
func (service *userInfoService) CreateUser(ctx context.Context, user User) (int64, error) {
	var err error = nil
	var params []string
	var resId int64
	if util.IsEmptyString(user.Name) {
		params = append(params, "name")
	}

	if util.IsEmptyString(user.PwdHash) {
		params = append(params, "password")
	}

	if v := len(params) > 0; v { // If one or more parameters were empty, creates a single error for all of them
		err = NewArgumentsRequiredError(params...)
	} else if age := user.Age; age < 1 || age > 150 {
		err = NewInvalidArgumentsError("age", "between 1 and 150")
	} else {
		resId, err = service.repository.InsertUser(user)
	}

	return resId, err
}

func (service *userInfoService) GetUser(ctx context.Context, id int64) (User, error) {
	if id < 1 {
		return User{}, NewInvalidArgumentsError("id", "must be 1 or greater")
	} else {
		return service.repository.GetUser(id)
	}
}
