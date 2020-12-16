package types

import (
	"context"
	"encoding/json"
)

// Writer is an object with the method IsAllowed.
type Writer interface {
	// IsAllowed returns an error, if the given user does not have the required
	// permission for the object. If it is allowed, it can also optionaly return
	// additional data as first return parameter.
	IsAllowed(ctx context.Context, userID int, payload map[string]json.RawMessage) (map[string]interface{}, error)
}

// WriterFunc is a function with the IsAllowed signature.
type WriterFunc func(ctx context.Context, userID int, payload map[string]json.RawMessage) (map[string]interface{}, error)

// IsAllowed calls the function.
func (f WriterFunc) IsAllowed(ctx context.Context, userID int, payload map[string]json.RawMessage) (map[string]interface{}, error) {
	return f(ctx, userID, payload)
}

// Reader is an object with a method to restrict fqfields.
type Reader interface {
	RestrictFQFields(ctx context.Context, userID int, fqfields []string, result map[string]bool) error
}

// ReaderFunc is a function with the IsAllowed signature.
type ReaderFunc func(ctx context.Context, userID int, fqfields []string, result map[string]bool) error

// RestrictFQFields calls the function.
func (f ReaderFunc) RestrictFQFields(ctx context.Context, userID int, fqfields []string, result map[string]bool) error {
	return f(ctx, userID, fqfields, result)
}
