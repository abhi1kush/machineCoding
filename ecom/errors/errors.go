package errors

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("NotFound")
var ErrUnintializedInstance = errors.New("not initialized")
var ErrSqlNOtFound = sql.ErrNoRows
