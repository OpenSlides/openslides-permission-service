package perm

import (
	"fmt"
)

// NotAllowedf is an error that sends a message to the client that indicates,
// that the user has not the required permissions.
func NotAllowedf(format string, a ...interface{}) error {
	return NotAllowedError{fmt.Sprintf(format, a...)}
}

// NotAllowedError tells, that the user does not have the required permission.
type NotAllowedError struct {
	reason string
}

func (e NotAllowedError) Error() string {
	return e.reason
}
