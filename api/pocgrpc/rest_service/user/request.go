package user

import (
	svcerr "github.com/EstebanFallaGlobant/globant-golang-bootcamp/api/pocgrpc/rest_service/error"
	"github.com/EstebanFallaGlobant/globant-golang-bootcamp/util"
)

type getUserRequest struct {
	authTokent string
	userID     int64
}

func (r getUserRequest) Validate() error {

	if util.IsEmptyString(r.authTokent) {
		return svcerr.NewInvalidRequestError(invRqstEmptyTknMsg)
	}

	if r.userID < 0 {
		return svcerr.NewInvalidRequestError(invRqstIDLessThanOne)
	}

	return nil
}
