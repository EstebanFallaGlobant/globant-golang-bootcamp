package entities

import (
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/error"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/util"
)

const (
	paramIDStr   = "ID"
	paramNameStr = "name"
	paramAgeStr  = "age"
	paramPassStr = "password"
	ruleMsgAge   = "must be between 1 and 150"
	ruleMsgID    = "must be 0 or greater"
)

type User struct {
	ID       int64
	Name     string
	PwdHash  string
	Age      uint8
	ParentID int64
}

type InitializationOption func(user *User) error

func NewUser(name, pwd string, age uint8, parentId int64, options ...InitializationOption) (User, error) {
	resultUser := User{
		Name:     name,
		PwdHash:  pwd,
		Age:      age,
		ParentID: parentId,
	}

	for _, option := range options {
		if err := option(&resultUser); err != nil {
			return User{}, err
		}

	}

	return resultUser, nil
}

func (u User) Validate() error {
	var params []string
	if util.IsEmptyString(u.Name) {
		params = append(params, paramNameStr)
	}

	if util.IsEmptyString(u.PwdHash) {
		params = append(params, paramPassStr)
	}

	if v := len(params) > 0; v { // If one or more parameters were empty, creates a single error for all of them
		return svcerr.NewArgumentsRequiredError(params...)
	}

	if age := u.Age; age < 1 || age > 150 {
		return svcerr.NewInvalidArgumentsError(paramAgeStr, ruleMsgAge)
	}

	if u.ID < 0 {
		return svcerr.NewInvalidArgumentsError(paramIDStr, ruleMsgID)
	}

	return nil
}
