package client

import "time"

const (
	genericToken     = "sometoken"
	genericID        = int64(2)
	genericPassword  = "testpassword"
	genericUsrName   = "Test User"
	genericUsrAge    = uint8(50)
	genericSleepTime = (1 * time.Second) / 2
)

const (
	authTokenCtxKey = "authToken"
)

const (
	paramAuthTokenName = "auth token"
)

const (
	ruleAuthTokenInvalidType = "be a valid generated token"
)
