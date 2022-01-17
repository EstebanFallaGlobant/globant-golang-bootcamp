package user

// Status strings used on error logs
const (
	errStatusKey  = "error"
	failStatusKey = "failed"
	nrmStatusKey  = "status"
)

// messages used on errors returned or on error logs
const (
	msgRqstNotParsed = "request could not be parsed"
	msgRspNotParsed  = "response could not be parsed"
	msgRqstDecoded   = "request decoded"
	msgRspEncoded    = "response encoded"
	msgRqstInvalid   = "invalid request"
)

// rule massages used on invalid argument errors or on error logs
const (
	ruleMsgCrtRqst = "must be 0"
	ruleMsgGetRqst = "must be 1 or greater"
	ruleMsgID      = "must be 1 or greater"
)

// name of parameters used on invalid argument errors or on error logs
const (
	paramIDStr = "ID"
)
