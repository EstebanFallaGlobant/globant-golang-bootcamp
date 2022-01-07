package user

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

const key = "Query Status"

type sqlRepository struct {
	logger log.Logger
	db     *sql.DB
}

func NewsqlRepository(logger log.Logger, db *sql.DB) *sqlRepository {
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &sqlRepository{
		logger: logger,
		db:     db,
	}
}

func (repo *sqlRepository) InsertUser(user User) (int64, error) {
	level.Info(repo.logger).Log(key, fmt.Sprintf("inserting user: %v", user))

	// Saves the name in lowercase to prevent user duplication on database
	result, err := repo.db.Exec(InsertUserQuery, strings.ToLower(user.Name), user.PwdHash, user.Age, user.parent)

	if err != nil {
		level.Error(repo.logger).Log(key, err)
		return 0, err
	} else {
		id, err := result.LastInsertId()

		if err != nil {
			level.Error(repo.logger).Log(key, err)
			return 0, err
		} else {
			level.Info(repo.logger).Log(key, fmt.Sprintf("user inserted with id: %d", id))
			return id, nil
		}
	}
}

func (repo *sqlRepository) GetUser(id int64) (User, error) {
	level.Info(repo.logger).Log(key, fmt.Sprintf("querying for user with id: %d", id))
	var user User

	if err := repo.db.QueryRow(SelectUserQuery, id).Scan(&user.Id, &user.PwdHash, &user.Name, &user.Age, &user.parent); err != nil {
		level.Error(repo.logger).Log(key, err)
		return User{}, err
	} else {
		level.Info(repo.logger).Log(key, fmt.Sprintf("user found: %v", user))
		user.Parent = user.parent.Int64
		return user, nil
	}
}

func (repo *sqlRepository) Close() {
	repo.db.Close()
}
