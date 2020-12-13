// Package permission provides tells, if a user has the permission to see or
// write an object.
package permission

import (
	"context"
	"fmt"
	"strings"

	"github.com/OpenSlides/openslides-permission-service/internal/definitions"
)

// Permission impelements the permission.Permission interface.
type Permission struct {
	coll map[string]Collection
}

// New returns a new permission service.
func New(dp DataProvider, os ...Option) *Permission {
	p := new(Permission)

	for _, o := range os {
		o(p)
	}

	if p.coll == nil {
		p.coll = openSlidesCollections(dp)
	}

	return p
}

// IsAllowed tells, if something is allowed.
func (ps *Permission) IsAllowed(ctx context.Context, name string, userID int, dataList []definitions.FqfieldData) ([]definitions.Addition, error) {
	collName := strings.SplitN(name, ".", 2)[0]
	coll, ok := ps.coll[collName]
	if !ok {
		return nil, clientError{fmt.Sprintf("unknown collection: `%s`", collName)}
	}

	additions := make([]definitions.Addition, len(dataList))
	for i, data := range dataList {
		addition, err := coll.IsAllowed(ctx, name, userID, data)
		if err != nil {
			return nil, isAllowedError{name: name, index: i, err: err}
		}

		additions[i] = addition
	}

	return additions, nil
}

// RestrictFQFields does currently nothing.
func (ps Permission) RestrictFQFields(ctx context.Context, userID int, fqfields []string) (map[string]bool, error) {
	restricted := make(map[definitions.Fqid]bool, len(fqfields))
	for _, fqfield := range fqfields {
		restricted[fqfield] = true
	}
	return restricted, nil
}

// AdditionalUpdate TODO
func (ps *Permission) AdditionalUpdate(ctx context.Context, updated definitions.FqfieldData) ([]definitions.ID, error) {
	return []definitions.ID{}, nil
}
