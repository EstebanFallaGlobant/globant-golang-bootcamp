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
	logger       log.Logger
	db           *sql.DB
	errorHandler sqlErrorHandler
}

type sqlErrorHandler interface {
	CreateUserServiceError(err error, user User) error
}

func NewsqlRepository(logger log.Logger, db *sql.DB, errorHandler sqlErrorHandler) *sqlRepository {
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &sqlRepository{
		logger:       logger,
		db:           db,
		errorHandler: errorHandler,
	}
}

func (repo *sqlRepository) InsertUser(user User) (int64, error) {
	level.Info(repo.logger).Log(key, fmt.Sprintf("inserting user: %v", user))

	statement, err := repo.db.Prepare(InsertUserQuery)
	if err != nil {
		level.Error(repo.logger).Log("sql preparation failed", err)
		return 0, err
	}

	// Saves the name in lowercase to prevent user duplication on database
	result, err := statement.Exec(strings.ToLower(user.Name), user.PwdHash, user.Age, user.parent)
	if err != nil {
		level.Error(repo.logger).Log(key, err)
		return 0, repo.errorHandler.CreateUserServiceError(err, user)
	} else {
		id, err := result.LastInsertId()

		if err != nil {
			level.Error(repo.logger).Log(key, err)
			return 0, repo.errorHandler.CreateUserServiceError(err, user)
		} else {
			level.Info(repo.logger).Log(key, fmt.Sprintf("user inserted with id: %d", id))
			return id, nil
		}
	}
}

func (repo *sqlRepository) GetUser(id int64) (User, error) {
	level.Info(repo.logger).Log(key, fmt.Sprintf("querying for user with id: %d", id))
	var user User
	{
		user = User{Id: id}
	}

	statement, err := repo.db.Prepare(SelectUserQuery)
	if err != nil {
		level.Error(repo.logger).Log("sql preparation failed", err)
		return User{}, repo.errorHandler.CreateUserServiceError(err, user)
	}

	if err := statement.QueryRow(id).Scan(&user.Id, &user.PwdHash, &user.Name, &user.Age, &user.parent); err != nil {
		level.Error(repo.logger).Log(key, err)
		return User{}, repo.errorHandler.CreateUserServiceError(err, user)
	} else {
		level.Info(repo.logger).Log(key, fmt.Sprintf("user found: %v", user))
		user.Parent = user.parent.Int64
		return user, nil
	}
}

func (repo *sqlRepository) Close() {
	repo.db.Close()
}
