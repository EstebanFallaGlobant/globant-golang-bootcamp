package user

const (
	InsertUserQuery = `INSERT INTO user_data(name,pwd_hash,age,parent_id) VALUES(?,?,?,?)`
	SelectUserQuery = `SELECT id, pwd_hash, name, age, parent_id FROM user_data WHERE id = ?`
)
