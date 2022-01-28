package user

import (
	"context"
	"errors"
	"testing"

	"os"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/entities"
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/error"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
)

func Test_WrapEndpoint(t *testing.T) {
	testCases := []struct {
		name           string
		middleWares    []endpoint.Middleware
		request        interface{}
		checkResponses func(t *testing.T, response interface{}, err error)
	}{
		{
			name: "Test single middleware wrapping",
			middleWares: []endpoint.Middleware{
				adderMiddleware,
			},
			request: 0,
			checkResponses: func(t *testing.T, response interface{}, err error) {
				assert.NoError(t, err)

				counter, ok := response.(int)

				if !ok {
					t.Fatal()
				}

				assert.EqualValues(t, 1, counter)
			},
		},
		{
			name: "Test multiple middleware wrapping",
			middleWares: []endpoint.Middleware{
				adderMiddleware,
				adderMiddleware,
				adderMiddleware,
			},
			request: 0,
			checkResponses: func(t *testing.T, response interface{}, err error) {
				assert.NoError(t, err)

				counter, ok := response.(int)

				if !ok {
					t.Fatal()
				}

				assert.EqualValues(t, 3, counter)
			},
		},
		{
			name: "Test middleware interrupt request chain",
			middleWares: []endpoint.Middleware{
				adderMiddleware,
			},
			request: "",
			checkResponses: func(t *testing.T, response interface{}, err error) {
				expectedError := svcerr.NewInvalidRequestError(invRqstNonParsed)
				assert.Nil(t, response)
				assert.ErrorAs(t, err, &expectedError)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			endpoint := wrapEndpoint(testEndpoint, tc.middleWares)

			response, err := endpoint(context.Background(), tc.request)

			tc.checkResponses(t, response, err)
		})
	}
}

func Test_MakeGetUserEndpoint(t *testing.T) {
	testCases := []struct {
		name          string
		searchID      int64
		serviceUser   entities.User
		serviceError  error
		request       interface{}
		checkResponse func(t *testing.T, response, expectedResponse interface{}, err error)
	}{
		{
			name:     "Test successful request",
			searchID: genericID,
			serviceUser: entities.User{
				ID:       genericID,
				Name:     genericUserName,
				Password: genericUserPassword,
				Age:      genericUserAge,
			},
			request: getUserRequest{
				authTokent: genericToken,
				userID:     genericID,
			},
			checkResponse: func(t *testing.T, response, expectedResponse interface{}, err error) {
				resp, ok := response.(getUserResponse)

				if !ok {
					t.Fatal(invRqstNonParsed)
				}

				assert.NoError(t, err)
				assert.EqualValues(t, expectedResponse, resp.user)
			},
		},
		{
			name:     "Test invalid request, non paseable",
			searchID: genericID,
			request:  genericID,
			checkResponse: func(t *testing.T, response, expectedResponse interface{}, err error) {
				expectedError := svcerr.NewInvalidRequestError(invRqstNonParsed)

				assert.Nil(t, response)
				assert.Error(t, err)
				assert.ErrorAs(t, err, &expectedError)
			},
		},
		{
			name:     "Test invalid request",
			searchID: genericID,
			request: getUserRequest{
				authTokent: genericToken,
				userID:     -1,
			},
			checkResponse: func(t *testing.T, response, expectedResponse interface{}, err error) {
				expectedError := svcerr.NewInvalidRequestError(invRqstIDLessThanOne)

				assert.Nil(t, response)
				assert.Error(t, err)
				assert.ErrorAs(t, err, &expectedError)
				assert.EqualError(t, err, expectedError.Error())
			},
		},
		{
			name:     "Test service returns error",
			searchID: genericID,
			request: getUserRequest{
				authTokent: genericToken,
				userID:     genericID,
			},
			serviceError: errors.New(genericErrorMsg),
			checkResponse: func(t *testing.T, response, expectedResponse interface{}, err error) {
				assert.Nil(t, response)
				assert.EqualError(t, err, genericErrorMsg)
			},
		},
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "timestamp", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := new(mockSvc)
			svc.On("GetUser", tc.searchID).
				Return(tc.serviceUser, tc.serviceError).
				Maybe()

			baseEndpoint := getUserEndpoint{
				svc:    svc,
				logger: logger,
			}
			endpoint := baseEndpoint.MakeGetUserEndpoint()

			response, err := endpoint(context.Background(), tc.request)

			tc.checkResponse(t, response, tc.serviceUser, err)
		})
	}
}

func Test_AddMiddlewares(t *testing.T) {
	testCases := []struct {
		name           string
		middlewares    []endpoint.Middleware
		targetEndpoint TargetEndpoint
		expectedError  error
		expectedCount  int
	}{
		{
			name:           "Test add single middleware",
			middlewares:    []endpoint.Middleware{adderMiddleware},
			targetEndpoint: GetUser,
			expectedCount:  1,
		},
		{
			name:           "Test add multiple middlewares",
			middlewares:    []endpoint.Middleware{adderMiddleware, adderMiddleware, adderMiddleware, adderMiddleware},
			targetEndpoint: GetUser,
			expectedCount:  4,
		},
		{
			name:           "Test undefined target error",
			middlewares:    []endpoint.Middleware{adderMiddleware},
			targetEndpoint: Undefined,
			expectedError:  svcerr.NewInvalidArgumentError(endpointTargetParamName, invArgEndpointTarget),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := getUserEndpoint{}
			count := endpoint.countMiddlewares

			assert.Zero(t, count(tc.targetEndpoint))

			err := endpoint.AddMiddlewares(tc.targetEndpoint, tc.middlewares...)

			assert.Equal(t, tc.expectedError, err)
			assert.EqualValues(t, tc.expectedCount, count(tc.targetEndpoint))
		})
	}
}

func Test_CountMiddlewares(t *testing.T) {
	testCases := []struct {
		name           string
		expectedCount  uint
		targetEndpoint TargetEndpoint
		middlewares    []endpoint.Middleware
	}{
		{
			name:           "Empty target",
			expectedCount:  0,
			targetEndpoint: GetUser,
		},
		{
			name:           "Single middleware in target",
			expectedCount:  1,
			targetEndpoint: GetUser,
			middlewares:    []endpoint.Middleware{adderMiddleware},
		},
		{
			name:           "Multiple middleware in target",
			expectedCount:  2,
			targetEndpoint: GetUser,
			middlewares:    []endpoint.Middleware{adderMiddleware, adderMiddleware},
		},
		{
			name:           "Undefined target",
			expectedCount:  0,
			targetEndpoint: Undefined,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := getUserEndpoint{}

			endpoint.getUserMiddlewares = tc.middlewares

			assert.EqualValues(t, tc.expectedCount, endpoint.countMiddlewares(tc.targetEndpoint))
		})
	}
}

func testEndpoint(ctx context.Context, request interface{}) (interface{}, error) {
	count, ok := request.(int)

	if !ok {
		return nil, svcerr.NewInvalidRequestError(invRqstNonParsed)
	}

	return count, nil
}

func adderMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		count, ok := request.(int)

		if !ok {
			return nil, svcerr.NewInvalidRequestError(invRqstNonParsed)
		}

		count++

		return next(ctx, count)
	}
}
