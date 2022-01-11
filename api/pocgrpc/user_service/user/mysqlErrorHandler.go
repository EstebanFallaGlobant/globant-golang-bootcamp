package user

import (
	"database/sql"
	"errors"

	"github.com/VividCortex/mysqlerr"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-sql-driver/mysql"
)

type MySQLErrorHandler struct {
	logger log.Logger
}

const (
	sqlErrHandlerKey = "status"
)

func (handler MySQLErrorHandler) CreateUserServiceError(err error, user User) error {
	level.Info(handler.logger).Log(sqlErrHandlerKey)

	if err == sql.ErrNoRows {
		return NewUserNotFoundError(user.Name, user.Id)
	}

	if driverErr, ok := err.(*mysql.MySQLError); ok {
		switch driverErr.Number {
		case mysqlerr.ER_BINLOG_UNSAFE_INSERT_TWO_KEYS:
			return NewUserAlreadyExistError(user.Name, user.Id)
		case mysqlerr.ER_RPL_INFO_DATA_TOO_LONG:
			return NewInvalidArgumentsError("user", "data too long")
		case mysqlerr.ER_NO_REFERENCED_ROW:
			return userNotUpdatedError{
				userId:  user.Id,
				message: driverErr.Message,
			}
		case mysqlerr.ER_ROW_IS_REFERENCED:
			return userNotUpdatedError{
				userId:  user.Id,
				message: driverErr.Message,
			}
		case mysqlerr.ER_AUTO_INCREMENT_CONFLICT:
			return userNotUpdatedError{
				userId:  user.Id,
				message: driverErr.Message,
			}
		default:
			return errors.New("operation failed")
		}
	}
	level.Error(handler.logger).Log("not a known mysql error", err)
	return err
}
