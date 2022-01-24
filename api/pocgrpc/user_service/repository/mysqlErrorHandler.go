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

func (handler MySQLErrorHandler) CreateUserServiceError(err error, user entities.User) error {
	level.Info(handler.Logger).Log(nrmStatusKey, "handling mySQL error")

	if err == sql.ErrNoRows {
		return svcerr.NewUserNotFoundError(user.Name, user.ID)
	}

	driverErr, ok := err.(*mysql.MySQLError)
	if !ok {
		level.Error(handler.Logger).Log(msgMysqlErrUnknown, err)
		return err
	}

	return sqlErrorToServiceError(driverErr, user)
}

func sqlErrorToServiceError(err *mysql.MySQLError, user entities.User) error {
	switch err.Number {
	case mysqlerr.ER_BINLOG_UNSAFE_INSERT_TWO_KEYS:
		return svcerr.NewUserAlreadyExistError(user.Name, user.ID)
	case mysqlerr.ER_RPL_INFO_DATA_TOO_LONG:
		return svcerr.NewInvalidArgumentsError(paramUsrStr, "data too long")
	case mysqlerr.ER_NO_REFERENCED_ROW:
		return svcerr.NewUserNotUpdatedError(user.ID, err.Message)
	case mysqlerr.ER_ROW_IS_REFERENCED:
		return svcerr.NewUserNotUpdatedError(user.ID, err.Message)
	case mysqlerr.ER_AUTO_INCREMENT_CONFLICT:
		return svcerr.NewUserNotUpdatedError(user.ID, err.Message)
	default:
		return errors.New(msgOperationFail)
	}
}
