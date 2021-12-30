package user

import (
	"context"

	"github.com/go-kit/log"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/util"
)

type Repository interface {
	InsertUser(user User) (int64, error)
}
type service struct {
	repository Repository
	log        log.Logger
}

func NewService(repo Repository, logger log.Logger) *service {
	return &service{
		repository: repo,
		log:        logger,
	}
}

// Creates an user if the information provided is valid.
//
// The users name and password must not be empty strings.
func (service *service) CreateUser(ctx context.Context, user User) (int64, error) {
	var err error = nil
	var params []string
	var resId int64
	if util.IsEmptyString(user.Name) {
		params = append(params, "name")
	}

	if util.IsEmptyString(user.PwdHash) {
		params = append(params, "password")
	}

	if id, v := int64(0), len(params) > 0; v { // If one or more parameters were empty, creates a single error for all of them
		err = NewArgumentsRequiredError(params...)
	} else if age := user.Age; age < 1 || age > 150 {
		err = NewInvalidArgumentsError("age", "between 1 and 150")
	} else {
		id, err = service.repository.InsertUser(
			NewUser(
				user.Name,
				user.PwdHash,
				user.Age,
				user.Parent))

		resId = id
	}

	return resId, err
}
