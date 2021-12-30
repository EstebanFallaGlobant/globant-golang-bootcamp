package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-kit/log"
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
	repo.logger.Log(key, fmt.Sprintf("inserting user: %v", user))

	// Saves the name in lowercase to prevent user duplication on database
	result, err := repo.db.Exec(InsertUserQuery, strings.ToLower(user.Name), user.PwdHash, user.Age, user.parent)

	if err != nil {
		repo.logger.Log(err)
		return 0, err
	} else {
		id, err := result.LastInsertId()

		if err != nil {
			repo.logger.Log(err)
			return 0, err
		} else {
			return id, nil
		}
	}
}

func (repo *sqlRepository) Close() {
	repo.db.Close()
}
