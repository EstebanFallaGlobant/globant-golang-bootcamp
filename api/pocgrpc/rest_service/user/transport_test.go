package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/entities"
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/error"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSvc struct {
	mock.Mock
}

func (mock mockSvc) GetUser(id int64) (entities.User, error) {
	args := mock.Called(id)
	return args.Get(0).(entities.User), args.Error(1)
}

type mockEndpoints struct {
	// mock.Mock
	svc    userService
	logger kitlog.Logger
}

func (mock mockEndpoints) MakeGetUserEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(getUserRequest)

		if !ok {
			return nil, svcerr.NewInvalidRequestError("")
		}

		if err := req.Validate(); err != nil {
			return nil, err
		}

		return mock.svc.GetUser(req.userID)
	}
}

func Test_GetUser(t *testing.T) {
	testCases := []struct {
		name          string
		searchID      int64
		expectedUsr   entities.User
		expectedErr   error
		searchHeader  map[string]string
		requestMethod string
		checkResult   func(t *testing.T, result *http.Response, expectedUser entities.User)
	}{
		{
			name:          "Test valid userID and token",
			searchID:      genericID,
			requestMethod: http.MethodGet,
			expectedUsr: entities.User{
				ID:       genericID,
				Name:     "Test user",
				Password: "TestPassword",
				Age:      20,
			},
			searchHeader: map[string]string{authTknHeaderName: genericToken},
			checkResult: func(t *testing.T, result *http.Response, expectedUser entities.User) {
				var user entities.User
				body, err := ioutil.ReadAll(result.Body)

				json.Unmarshal(body, &user)

				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, result.StatusCode)
				assert.EqualValues(t, expectedUser, user)
			},
		},
		{
			name:          "Test invalid header key",
			searchID:      genericID,
			requestMethod: http.MethodGet,
			searchHeader:  map[string]string{"Some-Header": genericToken},
			checkResult: func(t *testing.T, result *http.Response, expectedUser entities.User) {
				assert.Equal(t, http.StatusNotFound, result.StatusCode)
			},
		},
		{
			name:          "Test with empty auth token string",
			searchID:      genericID,
			requestMethod: http.MethodGet,
			searchHeader:  map[string]string{authTknHeaderName: " "},
			checkResult: func(t *testing.T, result *http.Response, expectedUser entities.User) {
				assert.Equal(t, http.StatusNotFound, result.StatusCode)
			},
		},
		{
			name:          "Test with wrong request method",
			searchID:      genericID,
			requestMethod: http.MethodPost,
			searchHeader:  map[string]string{authTknHeaderName: genericToken},
			checkResult: func(t *testing.T, result *http.Response, expectedUser entities.User) {
				assert.Equal(t, http.StatusMethodNotAllowed, result.StatusCode)
			},
		},
	}

	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "timestamp", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			svc := new(mockSvc)
			svc.On("GetUser", tc.searchID).Return(tc.expectedUsr, tc.expectedErr).
				Maybe()

			endpoints := mockEndpoints{
				svc:    svc,
				logger: logger,
			}

			httpServer := NewHTTPServer(logger, endpoints)

			server := httptest.NewServer(httpServer)
			defer server.Close()

			request := createNewRequest(tc.requestMethod, fmt.Sprintf("%s%s/%d", server.URL, getUserPath, tc.searchID), nil, tc.searchHeader)

			resp, err := http.DefaultClient.Do(request)

			if err != nil {
				t.Fatal(err)
			}

			tc.checkResult(t, resp, tc.expectedUsr)

		})
	}
}

func createNewRequest(method, url string, body io.Reader, headers map[string]string) *http.Request {
	request, _ := http.NewRequest(method, url, body)
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	return request
}

func Test_getErrorCode(t *testing.T) {
	testCases := []struct {
		name         string
		err          error
		expectedCode int
	}{
		{
			name:         "Invalid argument error",
			err:          svcerr.NewInvalidArgumentError(genericArgumentName, genericArgumentRule),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid request error",
			err:          svcerr.NewInvalidRequestError(""),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Generic error",
			err:          errors.New(genericArgumentRule),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code := getErrorCode(tc.err)

			assert.EqualValues(t, tc.expectedCode, code)
		})
	}
}
