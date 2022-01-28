package user

import "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/entities"

type errorResponse struct {
	Message string `json:"msg"`
	Code    int    `json:"status"`
}

type getUserResponse struct {
	user entities.User
}
