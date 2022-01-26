package entities

import (
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/error"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/util"
)

type User struct {
	ID       int64  `json:"user_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Age      uint8  `json:"age"`
	ParentID int64  `json:"parent_id"`
}

func (usr *User) Validate() []error {
	var errors []error
	if util.IsEmptyString(usr.Name) {
		errors = append(errors, svcerr.NewInvalidArgumentError(paramNameStr, ruleEmptyStr))
	}

	if util.IsEmptyString(usr.Password) {
		errors = append(errors, svcerr.NewInvalidArgumentError(paramPasswordStr, ruleEmptyStr))
	}

	if usr.Age < 1 {
		errors = append(errors, svcerr.NewInvalidArgumentError(paramAgeStr, ruleLessThanOne))
	}

	if usr.Age > MaxAllowedAge {
		errors = append(errors, svcerr.NewInvalidArgumentError(paramAgeStr, ruleGreaterThanAllowedAge))
	}

	if usr.ID < 0 {
		errors = append(errors, svcerr.NewInvalidArgumentError(paramIDStr, ruleLessThanZero))
	}

	if usr.ParentID < 0 {
		errors = append(errors, svcerr.NewInvalidArgumentError(paramParentIDStr, ruleLessThanZero))
	}

	return errors
}
