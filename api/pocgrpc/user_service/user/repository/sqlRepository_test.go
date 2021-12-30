package repository

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	kitlog "github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
)

const (
	InsertUserQueryTest = "INSERT INTO user_data"
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
			userInfo: NewUser(
				"Test user",
				"7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				uint8(20),
				0),
			expectedId: 1,
			sqlResultData: func(mock sqlmock.Sqlmock, user User, id int64) sqlmock.Sqlmock {
				mock.ExpectExec(InsertUserQueryTest).
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
			userInfo: NewUser(
				"Test user",
				"7bcf9d89298f1bfae16fa02ed6b61908fd2fa8de45dd8e2153a3c47300765328",
				uint8(20),
				0),
			expectedId: 0,
			sqlResultData: func(mock sqlmock.Sqlmock, user User, id int64) sqlmock.Sqlmock {
				mock.ExpectExec(InsertUserQueryTest).
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var logger kitlog.Logger
			{
				logger = kitlog.NewLogfmtLogger(os.Stderr)
				logger = kitlog.With(logger, "timestamp", kitlog.DefaultTimestampUTC)
				logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
			}

			db, mock, err := sqlmock.New()

			assert.NoError(t, err)
			mock = tc.sqlResultData(mock, tc.userInfo, tc.expectedId)

			mysqlrepo := NewsqlRepository(logger, db)

			var id int64
			id, err = mysqlrepo.InsertUser(tc.userInfo)

			assert.NoError(t, mock.ExpectationsWereMet())
			tc.checkResult(t, id, tc.expectedId, err)
		})
	}
}
