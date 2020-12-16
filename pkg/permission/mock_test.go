package permission_test

import (
	"context"
	"encoding/json"

	"github.com/OpenSlides/openslides-permission-service/internal/allowed"
	"github.com/OpenSlides/openslides-permission-service/internal/types"
	"github.com/OpenSlides/openslides-permission-service/pkg/permission"
)

func fakeCollections() []permission.Collection {
	return []permission.Collection{
		collectionMock{},
	}
}

type collectionMock struct{}

func (c collectionMock) WriteHandler() map[string]types.Writer {
	return map[string]types.Writer{
		"dummy_allowed":     allowedMock(true),
		"dummy_not_allowed": allowedMock(false),
	}
}

func (c collectionMock) ReadHandler() map[string]types.Reader {
	return map[string]types.Reader{
		"dummy": allowedMock(false),
	}
}

type allowedMock bool

func (a allowedMock) IsAllowed(ctx context.Context, userID int, data map[string]json.RawMessage) (map[string]interface{}, error) {
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
