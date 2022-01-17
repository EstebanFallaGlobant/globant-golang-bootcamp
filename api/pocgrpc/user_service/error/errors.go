package error

import (
	"errors"
	"fmt"
	"strings"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/util"
)

type InvalidArgumentsError struct {
	arguments map[string]string
	err       error
}

type ArgumentsRequiredError struct {
	err error
}

type InvalidRequestError struct {
	errMsg string
}

type UserNotFoundError struct {
	userName string
	userId   int64
}

type UserAlreadyExistsError struct {
	userName string
	userId   int64
}

type UserNotUpdatedError struct {
	userId  int64
	message string
}

func NewInvalidArgumentsError(keyval ...string) InvalidArgumentsError {
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

	return InvalidArgumentsError{
		err:       errors.New("invalid arguments"),
		arguments: args,
	}
}

func NewArgumentsRequiredError(arguments ...string) ArgumentsRequiredError {
	return ArgumentsRequiredError{
		err: errors.New(getArgumentRequiredErrorMessage(arguments)),
	}
}

func NewInvalidRequestError(message string) InvalidRequestError {
	return InvalidRequestError{
		errMsg: message,
	}
}

func NewUserNotFoundError(name string, id int64) UserNotFoundError {
	return UserNotFoundError{
		userName: name,
		userId:   id,
	}
}

func NewUserAlreadyExistError(name string, id int64) UserAlreadyExistsError {
	return UserAlreadyExistsError{
		userName: name,
		userId:   id,
	}
}

func NewUserNotUpdatedError(userID int64, message string) UserNotUpdatedError {
	return UserNotUpdatedError{
		userId:  userID,
		message: message,
	}
}

func (err UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("an user with the name \"%s\" was found with the id \"%d\"", err.userName, err.userId)
}

func (err UserNotFoundError) Error() string {
	if !util.IsEmptyString(err.userName) {
		return fmt.Sprintf("user with name \"%s\" not found", err.userName)
	}

	return fmt.Sprintf("user with id %d not found", err.userId)
}

func (err InvalidArgumentsError) Error() string {
	return err.getInvalidArgumentsErrorMessage()
}

func (err ArgumentsRequiredError) Error() string {
	return fmt.Sprintf("%v", err.err)
}

func (err InvalidRequestError) Error() string {
	return fmt.Sprintf("invalid request: %v", err.errMsg)
}

func (err UserNotUpdatedError) Error() string {
	return fmt.Sprintf("user: %d not updated with message: %s", err.userId, err.message)
}

func getArgumentRequiredErrorMessage(arguments []string) string {
	if len(arguments) > 1 {
		return strings.Join(arguments, ", ") + " are required"
	}
	return arguments[0] + " is required"

}

func (err *InvalidArgumentsError) getInvalidArgumentsErrorMessage() string {

	var message string
	for key, value := range err.arguments {
		message = message + fmt.Sprintf(", Parameter: \"%s\" Rule: \"%s\"", key, value)
	}

	return message
}
