// Package permission provides tells, if a user has the permission to see or
// write an object.
package permission

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
	"github.com/OpenSlides/openslides-permission-service/internal/perm"
)

// Permission impelements the permission.Permission interface.
type Permission struct {
	hs *handlerStore

	dp dataprovider.DataProvider
}

// New returns a new permission service.
func New(dp DataProvider) *Permission {
	p := &Permission{
		hs: newHandlerStore(),
		dp: dataprovider.DataProvider{External: dp},
	}

	for _, con := range openSlidesCollections(p.dp) {
		con.Connect(p.hs)
	}

	return p
}

// IsAllowed tells, if something is allowed.
func (ps *Permission) IsAllowed(ctx context.Context, name string, userID int, dataList []map[string]json.RawMessage) (bool, error) {
	superUser, err := ps.dp.IsSuperuser(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("checking for superuser: %w", err)
	}
	if superUser {
		return true, nil
	}

	// TODO: after all handlers are implemented. Move this code above the superUser check.
	handler, ok := ps.hs.writeHandler[name]
	if !ok {
		return false, fmt.Errorf("unknown collection: `%s`", name)
	}

	for i, data := range dataList {
		allowed, err := handler.IsAllowed(ctx, userID, data)
		if err != nil {
			return false, fmt.Errorf("action %d: %w", i, err)
		}
		if !allowed {
			return false, nil
		}
	}

	return true, nil
}

// superUserFields handles fields that the superuser is not allowed to see.
//
// Returns true, if the normal normal restricters should be skiped.
func superUserFields(result map[string]bool, collection string, fqfields []perm.FQField) (skip bool) {
	if collection == "personal_note" {
		return false
	}

	for _, k := range fqfields {
		if k.Collection == "user" && k.Field == "password" {
			continue
		}
		result[k.String()] = true
	}
	return true
}

// RestrictFQFields tells, if the given user can see the fqfields.
func (ps Permission) RestrictFQFields(ctx context.Context, userID int, fqfields []string) (map[string]bool, error) {
	data := make(map[string]bool, len(fqfields))

	superUser, err := ps.dp.IsSuperuser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("checking for superuser: %w", err)
	}

	grouped, err := groupFQFields(fqfields)
	if err != nil {
		return nil, fmt.Errorf("grouping fqfields: %w", err)
	}

	for name, fqfields := range grouped {
		if superUser {
			if superUserFields(data, name, fqfields) {
				continue
			}
		}

		handler, ok := ps.hs.readHandler[name]
		if !ok {
			return nil, fmt.Errorf("unknown collection: `%s`", name)
		}

		if err := handler.RestrictFQFields(ctx, userID, fqfields, data); err != nil {
			return nil, fmt.Errorf("restrict for collection %s: %w", name, err)
		}
	}
	return data, nil
}

func groupFQFields(fqfields []string) (map[string][]perm.FQField, error) {
	grouped := make(map[string][]perm.FQField)
	for _, f := range fqfields {
		fqfield, err := perm.ParseFQField(f)
		if err != nil {
			return nil, fmt.Errorf("decoding fqfield: %w", err)
		}
		grouped[fqfield.Collection] = append(grouped[fqfield.Collection], fqfield)
	}
	return grouped, nil
}

// AllRoutes returns the names of all read and write routes.
func (ps *Permission) AllRoutes() (readRoutes []string, writeRoutes []string) {
	rr := make([]string, 0, len(ps.hs.readHandler))
	for k := range ps.hs.readHandler {
		rr = append(rr, k)
	}

	wr := make([]string, 0, len(ps.hs.writeHandler))
	for k := range ps.hs.writeHandler {
		wr = append(wr, k)
	}
	return rr, wr
}

// DataProvider is the connection to the datastore. It returns the data
// required by the permission service.
type DataProvider interface {
	// If a field does not exist, it is not returned.
	Get(ctx context.Context, fqfields ...string) ([]json.RawMessage, error)
}

type handlerStore struct {
	writeHandler map[string]perm.ActionChecker
	readHandler  map[string]perm.RestricterChecker
}

func newHandlerStore() *handlerStore {
	return &handlerStore{
		writeHandler: make(map[string]perm.ActionChecker),
		readHandler:  make(map[string]perm.RestricterChecker),
	}
}

func (hs *handlerStore) RegisterRestricter(name string, reader perm.RestricterChecker) {
	if _, ok := hs.readHandler[name]; ok {
		panic(fmt.Sprintf("Read handler with name `%s` allready exists", name))
	}
	hs.readHandler[name] = reader
}

func (hs *handlerStore) RegisterAction(name string, writer perm.ActionChecker) {
	if _, ok := hs.writeHandler[name]; ok {
		panic(fmt.Sprintf("Write handler with name `%s` allready exists", name))
	}
	hs.writeHandler[name] = writer
}
