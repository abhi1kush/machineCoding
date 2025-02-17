package errors

import "errors"

var ErrNotFound = errors.New("NotFound")
var ErrUnintializedInstance = errors.New("not initialized")
