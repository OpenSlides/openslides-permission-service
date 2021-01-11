package collection

import (
	"context"

	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
	"github.com/OpenSlides/openslides-permission-service/internal/perm"
)

// Resource handels the permissions of resource objects.
func Resource(dp dataprovider.DataProvider) perm.ConnecterFunc {
	r := &resource{dp: dp}
	return func(s perm.HandlerStore) {
		s.RegisterReadHandler("resource", perm.ReadCheckerFunc(r.read))
	}
}

type resource struct {
	dp dataprovider.DataProvider
}

func (r *resource) read(ctx context.Context, userID int, fqfields []perm.FQField, result map[string]bool) error {
	for _, field := range fqfields {
		result[field.String()] = true
	}
	return nil
}
