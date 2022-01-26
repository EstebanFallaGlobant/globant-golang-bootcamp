package error

import (
	"fmt"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/util"
)

const (
	unknownParamName = "unknown"
	unknownParamRule = "is invalid"
	invRqstGenMsg    = "request is invalid"
	invRqstSpcMsg    = "invalid request:"
)

type InvalidArgumentError struct {
	name string
	rule string
}

func (err InvalidArgumentError) Error() string {
	if util.IsEmptyString(err.rule) {
		return fmt.Sprintf("parameter \"%s\" %s", err.name, unknownParamRule)
	}
	return fmt.Sprintf("parameter \"%s\" invalid: \"%s\"", err.name, err.rule)
}

func NewInvalidArgumentError(argName string, argRule string) InvalidArgumentError {
	if util.IsEmptyString(argName) {
		argName = unknownParamName
	}

	return InvalidArgumentError{
		name: argName,
		rule: argRule,
	}
}

type InvalidRequestError struct {
	message string
}

func (err InvalidRequestError) Error() string {
	if util.IsEmptyString(err.message) {
		return invRqstGenMsg
	}
	return fmt.Sprintf("%s %s", invRqstSpcMsg, err.message)
}

func NewInvalidRequestError(msg string) InvalidRequestError {
	return InvalidRequestError{
		message: msg,
	}
}
