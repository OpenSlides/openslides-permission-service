package permission_test

import (
	"context"
	"encoding/json"

	"github.com/OpenSlides/openslides-permission-service/internal/allowed"
	"github.com/OpenSlides/openslides-permission-service/internal/types"
)

func fakeCollections() []types.Connecter {
	return []types.Connecter{
		collectionMock{},
	}
}

type collectionMock struct{}

func (c collectionMock) Connect(s types.HandlerStore) {
	s.RegisterWriteHandler("dummy_allowed", allowedMock(true))
	s.RegisterWriteHandler("dummy_not_allowed", allowedMock(false))

	s.RegisterReadHandler("dummy", allowedMock(false))
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
