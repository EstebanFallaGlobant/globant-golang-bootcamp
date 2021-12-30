package repository

const (
	InsertUserQuery = `INSERT INTO user_data(name,pwd_hash,age,parent_id) VALUES(?,?,?,?)`
)
