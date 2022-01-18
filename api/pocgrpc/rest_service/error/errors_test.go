package error

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	genericParamName = "parameter"
	genericParamRule = "rule"
)

func Test_NewInvalidArgumentError(t *testing.T) {
	testCases := []struct {
		name        string
		param       string
		rule        string
		checkResult func(t *testing.T, result invalidArgumentError)
	}{
		{
			name:  "Test error creation with single parameter",
			param: genericParamName,
			rule:  genericParamRule,
			checkResult: func(t *testing.T, result invalidArgumentError) {
				assert.Contains(t, result.Error(), genericParamName)
				assert.Contains(t, result.Error(), genericParamRule)
			},
		},
		{
			name:  "Test error creation with empty parameter name",
			param: " ",
			rule:  genericParamRule,
			checkResult: func(t *testing.T, result invalidArgumentError) {
				assert.Contains(t, result.Error(), unknownParamName)
				assert.Contains(t, result.Error(), genericParamRule)
			},
		},
		{
			name:  "Test error creation with empty parameter rule",
			param: genericParamName,
			checkResult: func(t *testing.T, result invalidArgumentError) {
				assert.Contains(t, result.Error(), genericParamName)
				assert.Contains(t, result.Error(), unknownParamRule)
			},
		},
		{
			name: "Test error creation with param's name and rule empty",
			checkResult: func(t *testing.T, result invalidArgumentError) {
				assert.Contains(t, result.Error(), unknownParamName)
				assert.Contains(t, result.Error(), unknownParamRule)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NewInvalidArgumentError(tc.param, tc.rule)
			tc.checkResult(t, result)
		})
	}
}
