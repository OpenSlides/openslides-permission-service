package collection

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
	"github.com/OpenSlides/openslides-permission-service/internal/types"
)

// Generic is a helper object to create a collection with usual functions.
type Generic struct {
	dp         dataprovider.DataProvider
	collection string
	managePerm string
	readPerm   string
}

// NewGeneric creates a generic collection.
func NewGeneric(dp dataprovider.DataProvider, collection string, readPerm, managePerm string) *Generic {
	return &Generic{
		dp:         dp,
		collection: collection,
		managePerm: managePerm,
	}
}

// Connect sets the generic routs to the given reader and writer.
func (g *Generic) Connect(s types.HandlerStore) {
	s.RegisterWriteHandler(g.collection+".create", types.WriterFunc(g.Create))
	s.RegisterWriteHandler(g.collection+".update", types.WriterFunc(g.Modify))
	s.RegisterWriteHandler(g.collection+".delete", types.WriterFunc(g.Modify))

	s.RegisterReadHandler(g.collection, g)
}

func (g *Generic) check(ctx context.Context, meetingID int, userID int, payload map[string]json.RawMessage) (map[string]interface{}, error) {
	superUser, err := g.dp.IsSuperuser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if superUser {
		return nil, nil
	}

	if err := EnsurePerms(ctx, g.dp, userID, meetingID, g.managePerm); err != nil {
		return nil, fmt.Errorf("ensure manage permission: %w", err)
	}
	return nil, nil
}

// Create checks for the mermission to create a new object.
func (g *Generic) Create(ctx context.Context, userID int, payload map[string]json.RawMessage) (map[string]interface{}, error) {
	meetingID, err := MettingIDFromPayload(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("getting meeting id for create action: %w", err)
	}

	return g.check(ctx, meetingID, userID, payload)
}

// Modify checks for the permissions to alter an existing object.
func (g *Generic) Modify(ctx context.Context, userID int, payload map[string]json.RawMessage) (map[string]interface{}, error) {
	id, err := modelID(payload)
	if err != nil {
		return nil, fmt.Errorf("getting model id from payload: %w", err)
	}

	fqid := fmt.Sprintf("%s/%d", g.collection, id)
	meetingID, err := g.dp.MeetingFromModel(ctx, fqid)
	if err != nil {
		return nil, fmt.Errorf("getting meeting id for model %s: %w", fqid, err)
	}

	return g.check(ctx, meetingID, userID, payload)
}

// MettingIDFromPayload returns the meeting_id from the payload.
//
// It expects, that a field with the name "meeting_id" is in the payload.
func MettingIDFromPayload(ctx context.Context, payload map[string]json.RawMessage) (int, error) {
	var id int
	if err := json.Unmarshal(payload["meeting_id"], &id); err != nil {
		return 0, fmt.Errorf("no valid meeting id: %w", err)
	}

	return id, nil
}

// RestrictFQFields tells, if the user has the permission to see the requested
// fields.
func (g *Generic) RestrictFQFields(ctx context.Context, userID int, fqfields []string, result map[string]bool) error {
	if len(fqfields) == 0 {
		return nil
	}

	parts := strings.Split(fqfields[0], "/")
	meetingID, err := g.dp.MeetingFromModel(ctx, g.collection+"/"+parts[1])
	if err != nil {
		return fmt.Errorf("getting meeting from model: %w", err)
	}

	if err := EnsurePerms(ctx, g.dp, userID, meetingID, g.readPerm); err != nil {
		return nil
	}

	for _, fqfield := range fqfields {
		result[fqfield] = true
	}
	return nil
}

func modelID(data map[string]json.RawMessage) (int, error) {
	var id int
	if err := json.Unmarshal(data["id"], &id); err != nil {
		return 0, fmt.Errorf("no valid meeting id: %w", err)
	}
	return id, nil
}
