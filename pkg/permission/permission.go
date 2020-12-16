// Package permission provides tells, if a user has the permission to see or
// write an object.
package permission

import (
	"context"
	"fmt"
	"strings"

	"github.com/OpenSlides/openslides-permission-service/internal/definitions"
	"github.com/OpenSlides/openslides-permission-service/internal/types"
)

// Permission impelements the permission.Permission interface.
type Permission struct {
	colls        []Collection
	writeHandler map[string]types.Writer
	readHandler  map[string]types.Reader
}

// New returns a new permission service.
func New(dp DataProvider, os ...Option) *Permission {
	p := new(Permission)

	for _, o := range os {
		o(p)
	}

	if p.colls == nil {
		p.colls = openSlidesCollections(dp)
	}

	p.writeHandler = make(map[string]types.Writer)
	for _, coll := range p.colls {
		for name, handler := range coll.WriteHandler() {
			p.writeHandler[name] = handler
		}

		for name, handler := range coll.ReadHandler() {
			p.readHandler[name] = handler
		}
	}

	return p
}

// IsAllowed tells, if something is allowed.
func (ps *Permission) IsAllowed(ctx context.Context, name string, userID int, dataList []definitions.FqfieldData) ([]definitions.Addition, error) {
	handler, ok := ps.writeHandler[name]
	if !ok {
		return nil, clientError{fmt.Sprintf("unknown collection: `%s`", name)}
	}

	additions := make([]definitions.Addition, len(dataList))
	for i, data := range dataList {
		addition, err := handler.IsAllowed(ctx, userID, data)
		if err != nil {
			return nil, isAllowedError{name: name, index: i, err: err}
		}

		additions[i] = addition
	}

	return additions, nil
}

// RestrictFQFields does currently nothing.
func (ps Permission) RestrictFQFields(ctx context.Context, userID int, fqfields []string) (map[string]bool, error) {
	grouped := groupFQFields(fqfields)

	data := make(map[definitions.Fqid]bool, len(fqfields))

	for name, fqfields := range grouped {
		handler, ok := ps.readHandler[name]
		if !ok {
			return nil, clientError{fmt.Sprintf("unknown collection: `%s`", name)}
		}

		if err := handler.RestrictFQFields(ctx, userID, fqfields, data); err != nil {
			return nil, fmt.Errorf("restrict for collection %s: %w", name, err)
		}
	}
	return data, nil
}

// AdditionalUpdate TODO
func (ps *Permission) AdditionalUpdate(ctx context.Context, updated definitions.FqfieldData) ([]definitions.ID, error) {
	return []definitions.ID{}, nil
}

func groupFQFields(fqfields []string) map[string][]string {
	grouped := make(map[string][]string)
	for _, fqfield := range fqfields {
		parts := strings.Split(fqfield, "/")
		grouped[parts[0]] = append(grouped[parts[0]], fqfield)
	}
	return grouped
}
