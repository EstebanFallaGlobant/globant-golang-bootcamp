package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/entities"
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/error"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/util"
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
	CreateUserServiceError(err error, user entities.User) error
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

func (repo *sqlRepository) InsertUser(user entities.User) (int64, error) {
	level.Info(repo.logger).Log(key, fmt.Sprintf("inserting user: %v", user))

	statement, err := repo.db.Prepare(insertUserQuery)
	if err != nil {
		level.Error(repo.logger).Log("sql preparation failed", err)
		return 0, err
	}

	// Saves the name in lowercase to prevent user duplication on database
	result, err := statement.Exec(strings.ToLower(user.Name), user.PwdHash, user.Age, toSqlNullInt64(user.ParentID))
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

func (repo *sqlRepository) GetUser(id int64) (entities.User, error) {
	level.Info(repo.logger).Log(key, fmt.Sprintf("querying for user with id: %d", id))

	var user entities.User
	{
		user = entities.User{ID: id}
	}
	var parentId sql.NullInt64

	statement, err := repo.db.Prepare(selectUserQuery)
	if err != nil {
		level.Error(repo.logger).Log("sql preparation failed", err)
		return entities.User{}, repo.errorHandler.CreateUserServiceError(err, user)
	}

	if err := statement.QueryRow(id).Scan(&user.ID, &user.PwdHash, &user.Name, &user.Age, &parentId); err != nil {
		level.Error(repo.logger).Log(key, err)
		return entities.User{}, repo.errorHandler.CreateUserServiceError(err, user)
	}

	level.Info(repo.logger).Log(key, fmt.Sprintf("user found: %v", user))
	user.ParentID = parentId.Int64
	return user, nil

}
func (repo *sqlRepository) GetUserByName(name string) (entities.User, error) {
	level.Info(repo.logger).Log("querying for user by name: ", name)

	if util.IsEmptyString(name) {
		level.Info(repo.logger).Log("name is empty", name)
		return entities.User{}, svcerr.NewArgumentsRequiredError("name")
	}

	var user entities.User
	var parentId sql.NullInt64

	statement, err := repo.db.Prepare(selectUserByNameQuery)

	if err != nil {
		level.Error(repo.logger).Log("sql preparation failed", err)
		return entities.User{}, repo.errorHandler.CreateUserServiceError(err, user)
	}

	if err := statement.QueryRow(strings.ToLower(name)).Scan(&user.ID, &user.PwdHash, &user.Name, &user.Age, &parentId); err != nil {
		return entities.User{}, repo.errorHandler.CreateUserServiceError(err, user)
	}

	user.ParentID = parentId.Int64

	return user, nil
}

func (repo *sqlRepository) Close() {
	repo.db.Close()
}

func toSqlNullInt64(value int64) sql.NullInt64 {
	if value <= 0 {
		return sql.NullInt64{Int64: value}
	}

	return sql.NullInt64{
		Valid: true,
		Int64: value,
	}
}