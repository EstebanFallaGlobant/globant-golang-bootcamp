package errors

import (
	"errors"
	"fmt"
	"strings"
)

type invalidArgumentsError struct {
	Err error
}

type argumentsRequiredError struct {
	Err error
}

func NewInvalidArgumentsError(pairs ...string) *invalidArgumentsError {
	return &invalidArgumentsError{
		Err: errors.New(getInvalidArgumentsErrorMessage(pairs)),
	}
}

func NewArgumentsRequiredError(arguments ...string) *argumentsRequiredError {
	return &argumentsRequiredError{
		Err: errors.New(getArgumentRequiredErrorMessage(arguments)),
	}
}

func (err *invalidArgumentsError) Error() string {
	return fmt.Sprintf("%v", err.Err)
}

func (err *argumentsRequiredError) Error() string {
	return fmt.Sprintf("%v", err.Err)
}

func getArgumentRequiredErrorMessage(arguments []string) string {
	if len(arguments) > 1 {
		return strings.Join(arguments, ", ") + " are required"
	} else {
		return arguments[0] + " is required"
	}
}

func getInvalidArgumentsErrorMessage(pairs []string) string {
	if len(pairs) < 2 {
		panic("number of elements invalid")
	}
	if len(pairs)%2 != 0 {
		panic("number of strings on the slice is not even")
	}

	var message string
	for i := 0; i < len(pairs)-1; i += 2 {
		key, value := pairs[i], pairs[i+1]

		message = message + fmt.Sprintf(", Parameter: \"%s\" Rule: \"%s\"", key, value)
	}

	return message
}
