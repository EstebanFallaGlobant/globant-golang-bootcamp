package repository

// strings used as main key on error logs
const (
	sqlPrpFailKey = "sql preparation failed"
	nrmStatusKey  = "query Status"
)

// messages used on error logs
const (
	msgQueryingForUser = "querying for user"
	msgOperationFail   = "operation failed"
	msgMysqlErrUnknown = "not a known mysql error"
)

// name of parameters used on invalid argument or required argument errors or on error logs
const (
	paramUsrStr  = "user"
	paramNameStr = "name"
)
