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
	connecters   []types.Connecter
	writeHandler map[string]types.Writer
	readHandler  map[string]types.Reader
}

// New returns a new permission service.
func New(dp DataProvider, os ...Option) *Permission {
	p := &Permission{
		writeHandler: make(map[string]types.Writer),
		readHandler:  make(map[string]types.Reader),
	}

	for _, o := range os {
		o(p)
	}

	if p.connecters == nil {
		p.connecters = openSlidesCollections(dp)
	}

	p.writeHandler = make(map[string]types.Writer)
	for _, con := range p.connecters {
		con.Connect(p)
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

// RestrictFQFields tells, if the given user can see the fqfields.
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

// RegisterReadHandler registers a reader.
func (ps *Permission) RegisterReadHandler(name string, reader types.Reader) {
	ps.readHandler[name] = reader
}

// RegisterWriteHandler registers a writer.
func (ps *Permission) RegisterWriteHandler(name string, writer types.Writer) {
	ps.writeHandler[name] = writer
}

func groupFQFields(fqfields []string) map[string][]string {
	grouped := make(map[string][]string)
	for _, fqfield := range fqfields {
		parts := strings.Split(fqfield, "/")
		grouped[parts[0]] = append(grouped[parts[0]], fqfield)
	}
	return grouped
}
