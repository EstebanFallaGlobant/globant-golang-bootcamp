package user

import (
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	kitlog "github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	InsertUserQueryTest = "INSERT INTO user_data"
	GetUserQueryTest    = "SELECT"
)

func Test_Repository_UserCreation(t *testing.T) {
	testCases := []struct {
		name          string
		userInfo      User
		expectedId    int64
		sqlResultData func(mock sqlmock.Sqlmock, user User, expectedId int64) sqlmock.Sqlmock
		checkResult   func(t *testing.T, result, expectedResult interface{}, err error)
	}{
		{
			name: "Test user creation succesful",
			userInfo: User{
				Name:    "Test user",
				PwdHash: "7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				Age:     uint8(20),
				Parent:  0},
			expectedId: 1,
			sqlResultData: func(mock sqlmock.Sqlmock, user User, id int64) sqlmock.Sqlmock {
				mock.ExpectPrepare(InsertUserQueryTest).
					ExpectExec().
					WithArgs(strings.ToLower(user.Name), user.PwdHash, int(user.Age), user.parent).
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
			userInfo: User{
				Name:    "Test user",
				PwdHash: "7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				Age:     uint8(20),
				Parent:  0},
			expectedId: 0,
			sqlResultData: func(mock sqlmock.Sqlmock, user User, id int64) sqlmock.Sqlmock {
				mock.ExpectPrepare(InsertUserQueryTest).
					ExpectExec().
					WithArgs(strings.ToLower(user.Name), user.PwdHash, int(user.Age), user.parent).
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
		expectedUser User
		configMock   func(mock sqlmock.Sqlmock, id int64, resultUser User) sqlmock.Sqlmock
		checkResult  func(t *testing.T, expectedUser, resultUser User, resultError error)
	}{
		{
			name: "Test succesful query without parent",
			expectedUser: User{
				Id:      1,
				Name:    "Test user",
				PwdHash: "Test password",
				Age:     35,
				Parent:  0,
				parent:  sql.NullInt64{Int64: 0, Valid: false}},
			configMock: func(mock sqlmock.Sqlmock, id int64, resultUser User) sqlmock.Sqlmock {
				mock.ExpectPrepare(GetUserQueryTest).
					ExpectQuery().
					WithArgs(id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "pwd_hash", "name", "age", "parent_id"}).
							AddRow(resultUser.Id, resultUser.PwdHash, resultUser.Name, resultUser.Age, resultUser.parent))
				return mock
			},
			checkResult: func(t *testing.T, expectedUser, resultUser User, resultError error) {
				assert.NoError(t, resultError)
				assert.EqualValues(t, false, resultUser.parent.Valid)
				assert.EqualValues(t, expectedUser, resultUser)
			},
		},
		{
			name: "Test succesful query with parent",
			expectedUser: User{
				Id:      2,
				Name:    "Test user",
				PwdHash: "Test password",
				Age:     35,
				Parent:  1,
				parent:  sql.NullInt64{Int64: 1, Valid: true}},
			configMock: func(mock sqlmock.Sqlmock, id int64, resultUser User) sqlmock.Sqlmock {
				mock.ExpectPrepare(GetUserQueryTest).
					ExpectQuery().
					WithArgs(id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "pwd_hash", "name", "age", "parent_id"}).
							AddRow(resultUser.Id, resultUser.PwdHash, resultUser.Name, resultUser.Age, resultUser.parent))
				return mock
			},
			checkResult: func(t *testing.T, expectedUser, resultUser User, resultError error) {
				assert.NoError(t, resultError)
				assert.EqualValues(t, true, resultUser.parent.Valid)
				assert.EqualValues(t, expectedUser, resultUser)
			},
		},
		{
			name:         "Test invalid query",
			expectedUser: User{},
			configMock: func(mock sqlmock.Sqlmock, id int64, resultUser User) sqlmock.Sqlmock {
				mock.ExpectPrepare(GetUserQueryTest).
					ExpectQuery().
					WithArgs(id).
					WillReturnError(sql.ErrNoRows)
				return mock
			},
			checkResult: func(t *testing.T, expectedUser, resultUser User, resultError error) {
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
			sqlMock = tc.configMock(sqlMock, tc.expectedUser.Id, tc.expectedUser)

			sqlErrorHandler := new(mockSQLErrorHandler)
			sqlErrorHandler.On("CreateUserServiceError", mock.Anything, tc.expectedUser).Return(errors.New("test error"))

			repository := NewsqlRepository(logger, db, sqlErrorHandler)

			user, err := repository.GetUser(tc.expectedUser.Id)

			assert.NoError(t, sqlMock.ExpectationsWereMet())
			tc.checkResult(t, tc.expectedUser, user, err)
		})
	}
}
