package collection

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
)

// Generic is a helper object to create a collection with usual functions.
type Generic struct {
	dp         dataprovider.DataProvider
	collection string
	managePerm string
}

// CreateGeneric creates a generic collection.
func CreateGeneric(dp dataprovider.DataProvider, collection string, managePerm string) *Generic {
	return &Generic{
		dp:         dp,
		collection: collection,
		managePerm: managePerm,
	}
}

// IsAllowed impelements the permission.Collection interface.
func (g *Generic) IsAllowed(ctx context.Context, name string, userID int, data map[string]json.RawMessage) (map[string]interface{}, error) {

	superUser, err := IsSuperuser(ctx, userID, g.dp)
	if err != nil {
		return nil, err
	}
	if superUser {
		return nil, nil
	}

	nameParts := strings.Split(name, ".")
	if len(nameParts) != 2 {
		return nil, fmt.Errorf("TODO wrong name")
	}

	var meetingID int
	switch nameParts[1] {
	case "create":
		meetingID, err = g.create(ctx, userID, data)
		if err != nil {
			return nil, fmt.Errorf("getting meeting id for create action: %w", err)
		}

	case "update":
		fallthrough
	case "delete":
		meetingID, err = g.update(ctx, g.dp, userID, data)
		if err != nil {
			return nil, fmt.Errorf("getting meeting id for update action: %w", err)
		}

		meetingID, err = g.update(ctx, g.dp, userID, data)
		if err != nil {
			return nil, fmt.Errorf("getting meeting id for delete action: %w", err)
		}

	default:
		return nil, fmt.Errorf("TODO unknown name")
	}

	exists, err := DoesModelExists(ctx, "meeting/"+strconv.Itoa(meetingID), g.dp)
	if err != nil {
		return nil, fmt.Errorf("checking for meeting existing: %w", err)
	}

	if !exists {
		return nil, NotAllowedf("The meeting with id %d does not exist", meetingID)
	}

	if err := EnsurePerms(ctx, g.dp, userID, meetingID, g.managePerm); err != nil {
		return nil, fmt.Errorf("ensure manage permission: %w", err)
	}

	return nil, nil
}

func (g *Generic) create(ctx context.Context, userID int, data map[string]json.RawMessage) (int, error) {
	meetingID, err := meetingID(data)
	if err != nil {
		return 0, fmt.Errorf("getting meetingID: %w", err)
	}

	return meetingID, nil
}

func (g *Generic) update(ctx context.Context, dp dataprovider.DataProvider, userID int, data map[string]json.RawMessage) (int, error) {
	id, err := modelID(data)
	if err != nil {
		return 0, fmt.Errorf("getting model id: %w", err)
	}

	exists, err := DoesModelExists(ctx, fmt.Sprintf("%s/%d", g.collection, id), dp)
	if err != nil {
		return 0, fmt.Errorf("check that models does exist: %w", err)
	}
	if !exists {
		return 0, NotAllowedf("The %s with id %d does not exist", g.collection, id)
	}

	meetingID, err := MeetingFromModel(ctx, fmt.Sprintf("%s/%d", g.collection, id), dp)
	if err != nil {
		return 0, fmt.Errorf("getting meeting from model: %w", err)
	}
	return meetingID, nil
}

func meetingID(data map[string]json.RawMessage) (int, error) {
	var id int
	if err := json.Unmarshal(data["meeting_id"], &id); err != nil {
		return 0, fmt.Errorf("no valid meeting id: %w", err)
	}
	return id, nil
}

func modelID(data map[string]json.RawMessage) (int, error) {
	var id int
	if err := json.Unmarshal(data["id"], &id); err != nil {
		return 0, fmt.Errorf("no valid meeting id: %w", err)
	}
	return id, nil
}
