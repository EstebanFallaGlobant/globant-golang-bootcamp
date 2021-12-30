package user

import (
	"errors"
	"fmt"
	"strings"
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

func NewArgumentsRequiredError(arguments ...string) *argumentsRequiredError {
	return &argumentsRequiredError{
		Err: errors.New(getArgumentRequiredErrorMessage(arguments)),
	}
}

func NewInvalidRequestError() *invalidRequestError {
	return &invalidRequestError{
		Err: errors.New("request is invalid"),
	}
}

func (err *invalidArgumentsError) Error() string {
	return err.getInvalidArgumentsErrorMessage()
}

func (err *argumentsRequiredError) Error() string {
	return fmt.Sprintf("%v", err.Err)
}

func (err *invalidRequestError) Error() string {
	return fmt.Sprintf("%v", err.Err)
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
