package repository

const (
	insertUserQuery       = `INSERT INTO user_data(name,pwd_hash,age,parent_id) VALUES(?,?,?,?)`
	selectUserQuery       = `SELECT id, pwd_hash, name, age, parent_id FROM user_data WHERE id = ?`
	selectUserByNameQuery = `SELECT id, pwd_hash, name, age, parent_id FROM user_data WHERE name = ?`
)
