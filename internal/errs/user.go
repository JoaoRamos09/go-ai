package errs

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUserNotCreated = errors.New("user not created")
var ErrUserNotUpdated = errors.New("user not updated")
var ErrUserNotDeleted = errors.New("user not deleted")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrInvalidUserID = errors.New("invalid user id")
