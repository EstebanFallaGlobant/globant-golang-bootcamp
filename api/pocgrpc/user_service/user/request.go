package user

import (
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/entities"
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/error"
)

type createUserRequest struct {
	authToken string
	user      entities.User
}

type getUserRequest struct {
	authToken string
	id        int64
}

func (req createUserRequest) Validate() error {
	if req.user.ID != 0 {
		return svcerr.NewInvalidArgumentsError("ID", "must be 0")
	}

	return req.user.Validate()
}

func (req getUserRequest) Validate() error {
	if req.id < 1 {
		return svcerr.NewInvalidArgumentsError("ID", "must be 1 or greater")
	}

	return nil
}
