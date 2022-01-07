package user

import "database/sql"

type User struct {
	Id      int64
	Name    string
	PwdHash string
	Age     uint8
	Parent  int64
	parent  sql.NullInt64
}

type InitializationOption func(user *User) error

type createUserResponse struct {
	status error
	Id     int64
}

type createUserRequest struct {
	authToken string
	user      User
}

type getUserRequest struct {
	authToken string
	id        int64
}

type getUserResponse struct {
	status error
	user   User
}

func NewUser(name, pwd string, age uint8, parentId int64, options ...InitializationOption) (User, error) {
	resultUser := User{
		Name:    name,
		PwdHash: pwd,
		Age:     age,
		Parent:  parentId,
		parent: func(id int64) sql.NullInt64 {
			if id < 1 {
				return sql.NullInt64{
					Valid: false,
				}
			} else {
				return sql.NullInt64{
					Int64: id,
					Valid: true,
				}
			}
		}(parentId),
	}

	for _, option := range options {
		if err := option(&resultUser); err != nil {
			return User{}, err
		}

	}

	return resultUser, nil
}
