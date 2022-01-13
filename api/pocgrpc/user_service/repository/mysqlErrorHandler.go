package repository

import (
	"database/sql"
	"errors"

	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/entities"
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/user_service/error"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-sql-driver/mysql"
)

type MySQLErrorHandler struct {
	Logger log.Logger
}

const (
	sqlErrHandlerKey = "status"
)

func (handler MySQLErrorHandler) CreateUserServiceError(err error, user entities.User) error {
	level.Info(handler.Logger).Log(sqlErrHandlerKey, "handling mySQL error")

	if err == sql.ErrNoRows {
		return svcerr.NewUserNotFoundError(user.Name, user.ID)
	}

	if driverErr, ok := err.(*mysql.MySQLError); ok {
		switch driverErr.Number {
		case mysqlerr.ER_BINLOG_UNSAFE_INSERT_TWO_KEYS:
			return svcerr.NewUserAlreadyExistError(user.Name, user.ID)
		case mysqlerr.ER_RPL_INFO_DATA_TOO_LONG:
			return svcerr.NewInvalidArgumentsError("user", "data too long")
		case mysqlerr.ER_NO_REFERENCED_ROW:
			return svcerr.NewUserNotUpdatedError(user.ID, driverErr.Message)
		case mysqlerr.ER_ROW_IS_REFERENCED:
			return svcerr.NewUserNotUpdatedError(user.ID, driverErr.Message)
		case mysqlerr.ER_AUTO_INCREMENT_CONFLICT:
			return svcerr.NewUserNotUpdatedError(user.ID, driverErr.Message)
		default:
			return errors.New("operation failed")
		}
	}
	level.Error(handler.Logger).Log("not a known mysql error", err)
	return err
}
