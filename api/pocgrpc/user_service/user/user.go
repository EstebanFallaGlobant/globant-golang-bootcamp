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

func NewUser(name, pwd string, age uint8, parentId int64) User {
	return User{
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
}
