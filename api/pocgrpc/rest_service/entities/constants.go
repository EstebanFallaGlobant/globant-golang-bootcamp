package entities

const (
	MaxAllowedAge = uint8(150)
)

const (
	paramIDStr       = "id"
	paramNameStr     = "name"
	paramPasswordStr = "password"
	paramAgeStr      = "age"
	paramParentIDStr = "parent_id"
)

const (
	ruleEmptyStr              = "must not be empty or whitespace"
	ruleLessThanOne           = "must be greater or equal to 1"
	ruleGreaterThanAllowedAge = "must be less or equal to " + string(MaxAllowedAge)
	ruleLessThanZero          = "must be greater or equal to 0"
)

const (
	testPassword = "testpassword"
	testUsrName  = "Test User"
	testUsrAge   = uint8(50)
)
