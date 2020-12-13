package permission_test

import (
	"context"
	"encoding/json"

	"github.com/OpenSlides/openslides-permission-service/internal/allowed"
	"github.com/OpenSlides/openslides-permission-service/pkg/permission"
)

func fakeCollections() map[string]permission.Collection {
	return map[string]permission.Collection{
		"dummy_allowed":     allowedMock(true),
		"dummy_not_allowed": allowedMock(false),
	}
}

type allowedMock bool

func (a allowedMock) IsAllowed(ctx context.Context, name string, userID int, data map[string]json.RawMessage) (map[string]interface{}, error) {
	if !a {
		return nil, allowed.NotAllowed("Some reason here")
	}
	return nil, nil
}

func (a allowedMock) RestrictFQFields(ctx context.Context, userID int, fqfields []string, result map[string]bool) error {
	if !a {
		return nil
	}

	for _, fqfield := range fqfields {
		result[fqfield] = true
	}
	return nil
}
