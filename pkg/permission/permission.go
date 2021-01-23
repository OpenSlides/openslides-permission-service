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
	connecters   []perm.Connecter
	writeHandler map[string]perm.ActionChecker
	readHandler  map[string]perm.RestricterChecker

	dp dataprovider.DataProvider
}

// New returns a new permission service.
func New(dp DataProvider, os ...Option) *Permission {
	p := &Permission{
		writeHandler: make(map[string]perm.ActionChecker),
		readHandler:  make(map[string]perm.RestricterChecker),
		dp:           dataprovider.DataProvider{External: dp},
	}

	for _, o := range os {
		o(p)
	}

	if p.connecters == nil {
		p.connecters = openSlidesCollections(p.dp)
	}

	for _, con := range p.connecters {
		con.Connect(p)
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
	handler, ok := ps.writeHandler[name]
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

		handler, ok := ps.readHandler[name]
		if !ok {
			return nil, fmt.Errorf("unknown collection: `%s`", name)
		}

		if err := handler.RestrictFQFields(ctx, userID, fqfields, data); err != nil {
			return nil, fmt.Errorf("restrict for collection %s: %w", name, err)
		}
	}
	return data, nil
}

// RegisterRestricter registers a reader.
func (ps *Permission) RegisterRestricter(name string, reader perm.ReadChecker) {
	if _, ok := ps.readHandler[name]; ok {
		panic(fmt.Sprintf("Read handler with name `%s` allready exists", name))
	}
	ps.readHandler[name] = reader
}

// RegisterAction registers a writer.
func (ps *Permission) RegisterAction(name string, writer perm.WriteChecker) {
	if _, ok := ps.writeHandler[name]; ok {
		panic(fmt.Sprintf("Write handler with name `%s` allready exists", name))
	}
	ps.writeHandler[name] = writer
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
	rr := make([]string, 0, len(ps.readHandler))
	for k := range ps.readHandler {
		rr = append(rr, k)
	}
	wr := make([]string, 0, len(ps.writeHandler))
	for k := range ps.writeHandler {
		wr = append(wr, k)
	}
	return rr, wr
}
