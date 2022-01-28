package user

const (
	invRqstEmptyTknMsg   = "empty token"
	invRqstIDLessThanOne = "ID value less than one"
	invRqstNonParsed     = "request could not be parsed"
	invArgEndpointTarget = "must not be undefined"
)

const (
	genericToken        = "token"
	genericID           = int64(20)
	genericArgumentName = "argument"
	genericArgumentRule = "test argument rule"
	genericUserName     = "Test User"
	genericUserPassword = "testpassword"
	genericUserAge      = uint8(20)
	genericUserParent   = int64(1)
	genericErrorMsg     = "test error"
)

const (
	authTknHeaderName       = "Authentication-Token"
	userNotFoundKeyName     = "could not get user"
	endpointTargetParamName = "target"
)

const (
	getUserPath = "/user"
)
