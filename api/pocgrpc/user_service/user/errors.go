package user

import (
	"errors"
	"fmt"
	"strings"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/util"
)

type invalidArgumentsError struct {
	arguments map[string]string
	Err       error
}

type argumentsRequiredError struct {
	Err error
}

type invalidRequestError struct {
	Err error
}

type userNotFoundError struct {
	userName string
	userId   int64
}

type userAlreadyExistsError struct {
	userName string
	userId   int64
}

type userNotUpdatedError struct {
	userId  int64
	message string
}

func NewInvalidArgumentsError(keyval ...string) *invalidArgumentsError {
	if len(keyval) < 2 {
		panic("number of elements invalid")
	}
	if len(keyval)%2 != 0 {
		panic("number of strings on the slice is not even")
	}

	args := make(map[string]string, len(keyval)/2)

	for i := 0; i < len(keyval)-1; i += 2 {
		key, value := keyval[i], keyval[i+1]

		args[key] = value
	}
	return &invalidArgumentsError{
		Err:       errors.New("invalid arguments"),
		arguments: args,
	}
}

func NewArgumentsRequiredError(arguments ...string) argumentsRequiredError {
	return argumentsRequiredError{
		Err: errors.New(getArgumentRequiredErrorMessage(arguments)),
	}
}

func NewInvalidRequestError() invalidRequestError {
	return invalidRequestError{
		Err: errors.New("request is invalid"),
	}
}

func NewUserNotFoundError(name string, id int64) userNotFoundError {
	return userNotFoundError{
		userName: name,
		userId:   id,
	}
}

func NewUserAlreadyExistError(name string, id int64) userAlreadyExistsError {
	return userAlreadyExistsError{
		userName: name,
		userId:   id,
	}
}

func (err userAlreadyExistsError) Error() string {
	return fmt.Sprintf("an user with the name \"%s\" was found with the id \"%d\"", err.userName, err.userId)
}

func (err userNotFoundError) Error() string {
	if !util.IsEmptyString(err.userName) {
		return fmt.Sprintf("user with name \"%s\" not found", err.userName)
	} else {
		return fmt.Sprintf("user with id %d not found", err.userId)
	}
}

func (err invalidArgumentsError) Error() string {
	return err.getInvalidArgumentsErrorMessage()
}

func (err argumentsRequiredError) Error() string {
	return fmt.Sprintf("%v", err.Err)
}

func (err invalidRequestError) Error() string {
	return fmt.Sprintf("%v", err.Err)
}

func (err userNotUpdatedError) Error() string {
	return fmt.Sprintf("user: %d not updated with message: %s", err.userId, err.message)
}

func getArgumentRequiredErrorMessage(arguments []string) string {
	if len(arguments) > 1 {
		return strings.Join(arguments, ", ") + " are required"
	} else {
		return arguments[0] + " is required"
	}
}

func (err *invalidArgumentsError) getInvalidArgumentsErrorMessage() string {

	var message string
	for key, value := range err.arguments {
		message = message + fmt.Sprintf(", Parameter: \"%s\" Rule: \"%s\"", key, value)
	}

	return message
}
