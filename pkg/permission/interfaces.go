package permission

import (
	"context"
	"encoding/json"
)

// DataProvider is the connection to the datastore. It returns the data
// required by the permission service.
type DataProvider interface {
	// If a field does not exist, it is not returned.
	Get(ctx context.Context, fqfields ...string) ([]json.RawMessage, error)
}

// Collection represents one OpenSlides object like a motion or a User.
type Collection interface {
	IsAllowed(ctx context.Context, name string, userID int, data map[string]json.RawMessage) (map[string]interface{}, error)
	RestrictFQFields(ctx context.Context, userID int, fqfields []string, result map[string]bool) error
}
