package error

import (
	"fmt"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/util"
)

const (
	unknownParamName = "unknown"
	unknownParamRule = "is invalid"
)

type invalidArgumentError struct {
	name string
	rule string
}

func NewInvalidArgumentError(argName string, argRule string) invalidArgumentError {
	if util.IsEmptyString(argName) {
		argName = unknownParamName
	}

	return invalidArgumentError{
		name: argName,
		rule: argRule,
	}
}

func (err invalidArgumentError) Error() string {
	if util.IsEmptyString(err.rule) {
		return fmt.Sprintf("parameter \"%s\" %s", err.name, unknownParamRule)
	}
	return fmt.Sprintf("parameter \"%s\" invalid: \"%s\"", err.name, err.rule)
}
