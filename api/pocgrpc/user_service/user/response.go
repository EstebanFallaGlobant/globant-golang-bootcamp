package user

import "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/entities"

type getUserResponse struct {
	status error
	user   entities.User
}

type createUserResponse struct {
	status error
	Id     int64
}
