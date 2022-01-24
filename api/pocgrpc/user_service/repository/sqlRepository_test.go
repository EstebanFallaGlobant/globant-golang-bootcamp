package repository

import (
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/entities"
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/error"
	kitlog "github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test queries
const (
	insertUserQueryTest = "INSERT INTO user_data"
	getUserQueryTest    = "SELECT"
)

type mockSQLErrorHandler struct {
	mock.Mock
}

func (mock mockSQLErrorHandler) CreateUserServiceError(err error, user entities.User) error {
	args := mock.Called(err, user)
	return args.Error(0)
}

func Test_Repository_UserCreation(t *testing.T) {
	testCases := []struct {
		name          string
		userInfo      entities.User
		expectedId    int64
		sqlResultData func(mock sqlmock.Sqlmock, user entities.User, expectedId int64) sqlmock.Sqlmock
		checkResult   func(t *testing.T, result, expectedResult interface{}, err error)
	}{
		{
			name: "Test user creation succesful",
			userInfo: entities.User{
				Name:     "Test user",
				PwdHash:  "7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				Age:      uint8(20),
				ParentID: 0},
			expectedId: 1,
			sqlResultData: func(mock sqlmock.Sqlmock, user entities.User, id int64) sqlmock.Sqlmock {
				mock.ExpectPrepare(insertUserQueryTest).
					ExpectExec().
					WithArgs(strings.ToLower(user.Name), user.PwdHash, int(user.Age), toSqlNullInt64(user.ParentID)).
					WillReturnResult(sqlmock.NewResult(id, 1))
				return mock
			},
			checkResult: func(t *testing.T, result, expectedResult interface{}, err error) {

				id, idOk := result.(int64)
				expectedId, expOk := expectedResult.(int64)

				if !idOk || !expOk {
					t.Fatal(errors.New("The result or expected result couldn't be parsed as int64"))
				}

				assert.NoError(t, err)
				assert.EqualValues(t, expectedId, id)
			},
		},
		{
			name: "Test user creation failed",
			userInfo: entities.User{
				Name:     "Test user",
				PwdHash:  "7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				Age:      uint8(20),
				ParentID: 0},
			expectedId: 0,
			sqlResultData: func(mock sqlmock.Sqlmock, user entities.User, id int64) sqlmock.Sqlmock {
				mock.ExpectPrepare(insertUserQueryTest).
					ExpectExec().
					WithArgs(strings.ToLower(user.Name), user.PwdHash, int(user.Age), toSqlNullInt64(user.ParentID)).
					WillReturnError(errors.New("some sql error"))
				return mock
			},
			checkResult: func(t *testing.T, result, expectedResult interface{}, err error) {

				id, idOk := result.(int64)
				expectedId, expOk := expectedResult.(int64)

				if !idOk || !expOk {
					t.Fatal(errors.New("The result or expected result couldn't be parsed as int64"))
				}

				assert.Error(t, err)
				assert.EqualValues(t, expectedId, id)
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

			db, sqlmock, err := sqlmock.New()
			assert.NoError(t, err)
			sqlmock = tc.sqlResultData(sqlmock, tc.userInfo, tc.expectedId)

			sqlErrorHandler := new(mockSQLErrorHandler)
			sqlErrorHandler.On("CreateUserServiceError", mock.Anything, tc.userInfo).Return(errors.New("test error"))

			mysqlrepo := NewsqlRepository(logger, db, sqlErrorHandler)

			var id int64
			id, err = mysqlrepo.InsertUser(tc.userInfo)

			assert.NoError(t, sqlmock.ExpectationsWereMet())
			tc.checkResult(t, id, tc.expectedId, err)
		})
	}
}

func Test_Repository_GetUser(t *testing.T) {
	testCases := []struct {
		name         string
		expectedUser entities.User
		configMock   func(mock sqlmock.Sqlmock, id int64, resultUser entities.User) sqlmock.Sqlmock
		checkResult  func(t *testing.T, expectedUser, resultUser entities.User, resultError error)
	}{
		{
			name: "Test succesful query without parent",
			expectedUser: entities.User{
				ID:       1,
				Name:     "Test user",
				PwdHash:  "Test password",
				Age:      35,
				ParentID: 0,
			},
			configMock: func(mock sqlmock.Sqlmock, id int64, resultUser entities.User) sqlmock.Sqlmock {
				mock.ExpectPrepare(getUserQueryTest).
					ExpectQuery().
					WithArgs(id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "pwd_hash", "name", "age", "parent_id"}).
							AddRow(resultUser.ID, resultUser.PwdHash, resultUser.Name, resultUser.Age, toSqlNullInt64(resultUser.ParentID)))
				return mock
			},
			checkResult: func(t *testing.T, expectedUser, resultUser entities.User, resultError error) {
				assert.NoError(t, resultError)
				assert.EqualValues(t, expectedUser, resultUser)
			},
		},
		{
			name: "Test succesful query with parent",
			expectedUser: entities.User{
				ID:       2,
				Name:     "Test user",
				PwdHash:  "Test password",
				Age:      35,
				ParentID: 1,
			},
			configMock: func(mock sqlmock.Sqlmock, id int64, resultUser entities.User) sqlmock.Sqlmock {
				mock.ExpectPrepare(getUserQueryTest).
					ExpectQuery().
					WithArgs(id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "pwd_hash", "name", "age", "parent_id"}).
							AddRow(resultUser.ID, resultUser.PwdHash, resultUser.Name, resultUser.Age, toSqlNullInt64(resultUser.ParentID)))
				return mock
			},
			checkResult: func(t *testing.T, expectedUser, resultUser entities.User, resultError error) {
				assert.NoError(t, resultError)
				assert.EqualValues(t, expectedUser, resultUser)
			},
		},
		{
			name:         "Test invalid query",
			expectedUser: entities.User{},
			configMock: func(mock sqlmock.Sqlmock, id int64, resultUser entities.User) sqlmock.Sqlmock {
				mock.ExpectPrepare(getUserQueryTest).
					ExpectQuery().
					WithArgs(id).
					WillReturnError(sql.ErrNoRows)
				return mock
			},
			checkResult: func(t *testing.T, expectedUser, resultUser entities.User, resultError error) {
				assert.Error(t, resultError)
				assert.IsType(t, errors.New("test error"), resultError)
				assert.EqualValues(t, expectedUser, resultUser)
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
			db, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			sqlMock = tc.configMock(sqlMock, tc.expectedUser.ID, tc.expectedUser)

			sqlErrorHandler := new(mockSQLErrorHandler)
			sqlErrorHandler.On("CreateUserServiceError", mock.Anything, tc.expectedUser).Return(errors.New("test error"))

			repository := NewsqlRepository(logger, db, sqlErrorHandler)

			user, err := repository.GetUser(tc.expectedUser.ID)

			assert.NoError(t, sqlMock.ExpectationsWereMet())
			tc.checkResult(t, tc.expectedUser, user, err)
		})
	}
}

func Test_GetUserByName(t *testing.T) {
	testCases := []struct {
		name         string
		expectedUser entities.User
		configMock   func(mock sqlmock.Sqlmock, name string, resultUser entities.User) sqlmock.Sqlmock
		checkResult  func(t *testing.T, expectedUser, resultUser entities.User, resultError error)
	}{
		{
			name: "Test successful query",
			expectedUser: entities.User{
				ID:   1,
				Name: "Test User",
			},
			configMock: func(mock sqlmock.Sqlmock, name string, resultUser entities.User) sqlmock.Sqlmock {
				mock.ExpectPrepare(getUserQueryTest).
					ExpectQuery().
					WithArgs(strings.ToLower(name)).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "pwd_hash", "name", "age", "parent_id"}).
							AddRow(resultUser.ID, resultUser.PwdHash, resultUser.Name, resultUser.Age, toSqlNullInt64(resultUser.ParentID)))
				return mock
			},
			checkResult: func(t *testing.T, expectedUser, resultUser entities.User, resultError error) {
				assert.NoError(t, resultError)
				assert.EqualValues(t, expectedUser, resultUser)
			},
		},
		{
			name: "Test Upper case name must return Lower case user",
			expectedUser: entities.User{
				ID:   1,
				Name: "Test User",
			},
			configMock: func(mock sqlmock.Sqlmock, name string, resultUser entities.User) sqlmock.Sqlmock {
				mock.ExpectPrepare(getUserQueryTest).
					ExpectQuery().
					WithArgs(strings.ToLower(name)).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "pwd_hash", "name", "age", "parent_id"}).
							AddRow(resultUser.ID, resultUser.PwdHash, strings.ToLower(resultUser.Name), resultUser.Age, toSqlNullInt64(resultUser.ParentID)))
				return mock
			},
			checkResult: func(t *testing.T, expectedUser, resultUser entities.User, resultError error) {
				resultUser.Name = strings.ToLower(resultUser.Name)

				assert.NoError(t, resultError)
				assert.EqualValues(t, strings.ToLower(expectedUser.Name), resultUser.Name)
			},
		},
		{
			name: "Test with empty name",
			expectedUser: entities.User{
				Name: " ",
			},
			configMock: func(mock sqlmock.Sqlmock, name string, resultUser entities.User) sqlmock.Sqlmock {
				return mock
			},
			checkResult: func(t *testing.T, expectedUser, resultUser entities.User, resultError error) {
				expectedError := svcerr.NewArgumentsRequiredError("name")
				assert.Error(t, resultError)
				assert.ErrorAs(t, resultError, &expectedError)
				assert.EqualValues(t, entities.User{}, resultUser)
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
			db, sqlMock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			sqlMock = tc.configMock(sqlMock, tc.expectedUser.Name, tc.expectedUser)

			sqlErrorHandler := new(mockSQLErrorHandler)
			sqlErrorHandler.On("CreateUserServiceError", mock.Anything, tc.expectedUser).Return(errors.New("test error"))

			repository := NewsqlRepository(logger, db, sqlErrorHandler)

			user, err := repository.GetUserByName(tc.expectedUser.Name)

			assert.NoError(t, sqlMock.ExpectationsWereMet())
			tc.checkResult(t, tc.expectedUser, user, err)
		})
	}
}
