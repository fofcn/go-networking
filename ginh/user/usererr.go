package user

import "errors"

var (
	UserHasRegisteredErr = errors.New("User has registerd")
	PwdFmtIncorrect      = errors.New("passowrd format is incorrect")
)
