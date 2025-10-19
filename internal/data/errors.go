// Filename: internal/data/errors.go

package data

import (
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")
var ErrEditConflict = errors.New("edit Conflict")

var ErrCourseNotFound = errors.New("course not found")
var ErrPostingNotFound = errors.New("posting not found")
