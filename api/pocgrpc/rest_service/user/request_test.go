package user

import (
	"testing"

	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/error"
	"github.com/stretchr/testify/assert"
)

func Test_getUserRequest_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		testRequest getUserRequest
		checkResult func(t *testing.T, resultError error)
	}{
		{
			name: "Valid request",
			testRequest: getUserRequest{
				authTokent: genericToken,
				userID:     genericID,
			},
			checkResult: func(t *testing.T, resultError error) {
				assert.NoError(t, resultError)
			},
		},
		{
			name: "Empty token",
			testRequest: getUserRequest{
				userID: genericID,
			},
			checkResult: func(t *testing.T, resultError error) {
				expectedError := svcerr.NewInvalidRequestError(invRqstEmptyTknMsg)
				assert.Error(t, resultError)
				assert.ErrorAs(t, resultError, &expectedError)
			},
		},
		{
			name: "User ID -1",
			testRequest: getUserRequest{
				authTokent: genericToken,
				userID:     -1,
			},
			checkResult: func(t *testing.T, resultError error) {
				expectedError := svcerr.NewInvalidRequestError(invRqstIDLessThanOne)
				assert.Error(t, resultError)
				assert.ErrorAs(t, resultError, &expectedError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.checkResult(t, tc.testRequest.Validate())
		})
	}
}
