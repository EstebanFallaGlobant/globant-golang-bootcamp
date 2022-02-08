package transform

import (
	"testing"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/entities"
	pb "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/pb"
	"github.com/stretchr/testify/assert"
)

const (
	defaultName     = "Test Name"
	defaultPassword = "TestPassword"
	defaultAge      = uint8(20)
	defaultID       = int64(2)
	defaultParentID = int64(1)
)

func Test_FromUserToPbUser(t *testing.T) {
	testCases := []struct {
		name             string
		expectedName     string
		expectedPassword string
		expectedAge      uint8
		expectedParent   int64
		expectedID       int64
		expectedError    error
		checkResult      func(t *testing.T, resultUser *pb.User, expectedUser *entities.User)
	}{
		{
			name:             "Test user convertion with ID zero",
			expectedName:     defaultName,
			expectedAge:      defaultAge,
			expectedPassword: defaultPassword,
			checkResult: func(t *testing.T, resultUser *pb.User, expectedUser *entities.User) {
				assert.Zero(t, resultUser.Id)
				assert.Zero(t, resultUser.ParentId)
				assert.EqualValues(t, expectedUser.Age, resultUser.Age)
				assert.EqualValues(t, expectedUser.Name, resultUser.Name)
				assert.EqualValues(t, expectedUser.Password, resultUser.PwdHash)
			},
		},
		{
			name:             "Test user convertion with ID one",
			expectedID:       defaultID,
			expectedName:     defaultName,
			expectedAge:      defaultAge,
			expectedPassword: defaultPassword,
			checkResult: func(t *testing.T, resultUser *pb.User, expectedUser *entities.User) {
				assert.Zero(t, resultUser.ParentId)
				assert.EqualValues(t, expectedUser.ID, resultUser.Id)
				assert.EqualValues(t, expectedUser.Age, resultUser.Age)
				assert.EqualValues(t, expectedUser.Name, resultUser.Name)
				assert.EqualValues(t, expectedUser.Password, resultUser.PwdHash)
			},
		},
		{
			name:             "Test user convertion with ParentID one",
			expectedID:       defaultID,
			expectedName:     defaultName,
			expectedAge:      defaultAge,
			expectedPassword: defaultPassword,
			checkResult: func(t *testing.T, resultUser *pb.User, expectedUser *entities.User) {
				assert.Zero(t, resultUser.ParentId)
				assert.EqualValues(t, expectedUser.ID, resultUser.Id)
				assert.EqualValues(t, expectedUser.Age, resultUser.Age)
				assert.EqualValues(t, expectedUser.Name, resultUser.Name)
				assert.EqualValues(t, expectedUser.Password, resultUser.PwdHash)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := entities.User{
				ID:       tc.expectedID,
				Name:     tc.expectedName,
				Password: tc.expectedPassword,
				Age:      tc.expectedAge,
				ParentID: tc.expectedParent,
			}

			result := FromUserToPbUser(user)

			tc.checkResult(t, &result, &user)
		})
	}
}
