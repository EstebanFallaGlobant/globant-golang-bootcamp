package service

import "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/user/repository"

type Repository interface {
	InsertUser(user repository.User) (int64, error)
}
