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
		checkResult func(t *testing.T, result InvalidArgumentError)
	}{
		{
			name:  "Test error creation with single parameter",
			param: genericParamName,
			rule:  genericParamRule,
			checkResult: func(t *testing.T, result InvalidArgumentError) {
				assert.Contains(t, result.Error(), genericParamName)
				assert.Contains(t, result.Error(), genericParamRule)
			},
		},
		{
			name:  "Test error creation with empty parameter name",
			param: " ",
			rule:  genericParamRule,
			checkResult: func(t *testing.T, result InvalidArgumentError) {
				assert.Contains(t, result.Error(), unknownParamName)
				assert.Contains(t, result.Error(), genericParamRule)
			},
		},
		{
			name:  "Test error creation with empty parameter rule",
			param: genericParamName,
			checkResult: func(t *testing.T, result InvalidArgumentError) {
				assert.Contains(t, result.Error(), genericParamName)
				assert.Contains(t, result.Error(), unknownParamRule)
			},
		},
		{
			name: "Test error creation with param's name and rule empty",
			checkResult: func(t *testing.T, result InvalidArgumentError) {
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

func Test_NewInvalidRequestError(t *testing.T) {
	testCases := []struct {
		name        string
		message     string
		checkResult func(t *testing.T, resultError error)
	}{
		{
			name:    "Valid message",
			message: "some message",
			checkResult: func(t *testing.T, resultError error) {
				assert.Error(t, resultError)
				assert.Contains(t, resultError.Error(), invRqstSpcMsg)
				assert.Contains(t, resultError.Error(), "some message")
			},
		},
		{
			name:    "Empty message",
			message: " ",
			checkResult: func(t *testing.T, resultError error) {
				assert.Error(t, resultError)
				assert.Equal(t, resultError.Error(), invRqstGenMsg)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.checkResult(t, NewInvalidRequestError(tc.message))
		})
	}
}
