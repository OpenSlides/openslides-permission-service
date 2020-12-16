package permission

import (
	"context"
	"encoding/json"

	"github.com/OpenSlides/openslides-permission-service/internal/types"
)

// DataProvider is the connection to the datastore. It returns the data
// required by the permission service.
type DataProvider interface {
	// If a field does not exist, it is not returned.
	Get(ctx context.Context, fqfields ...string) ([]json.RawMessage, error)
}

// Collection represents one OpenSlides object like a motion or a User.
type Collection interface {
	WriteHandler() map[string]types.Writer
	ReadHandler() map[string]types.Reader
}
