package collection

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
	"github.com/OpenSlides/openslides-permission-service/internal/types"
)

// Create checks for the mermission to create a new object.
func Create(dp dataprovider.DataProvider, perm, collection string) types.Writer {
	return types.WriterFunc(func(ctx context.Context, userID int, payload map[string]json.RawMessage) (map[string]interface{}, error) {
		meetingID, err := MettingIDFromPayload(ctx, payload)
		if err != nil {
			return nil, fmt.Errorf("getting meeting id for create action: %w", err)
		}

		return check(ctx, dp, perm, meetingID, userID, payload)
	})
}

// Modify checks for the permissions to alter an existing object.
func Modify(dp dataprovider.DataProvider, perm, collection string) types.Writer {
	return types.WriterFunc(func(ctx context.Context, userID int, payload map[string]json.RawMessage) (map[string]interface{}, error) {
		id, err := modelID(payload)
		if err != nil {
			return nil, fmt.Errorf("getting model id from payload: %w", err)
		}

		fqid := fmt.Sprintf("%s/%d", collection, id)
		meetingID, err := dp.MeetingFromModel(ctx, fqid)
		if err != nil {
			return nil, fmt.Errorf("getting meeting id for model %s: %w", fqid, err)
		}

		return check(ctx, dp, perm, meetingID, userID, payload)
	})
}

func check(ctx context.Context, dp dataprovider.DataProvider, managePerm string, meetingID int, userID int, payload map[string]json.RawMessage) (map[string]interface{}, error) {
	superUser, err := dp.IsSuperuser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if superUser {
		return nil, nil
	}

	if err := EnsurePerms(ctx, dp, userID, meetingID, managePerm); err != nil {
		return nil, fmt.Errorf("ensure manage permission: %w", err)
	}
	return nil, nil
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

func modelID(data map[string]json.RawMessage) (int, error) {
	var id int
	if err := json.Unmarshal(data["id"], &id); err != nil {
		return 0, fmt.Errorf("no valid meeting id: %w", err)
	}
	return id, nil
}

// Restrict tells, if the user has the permission to see the requested
// fields.
func Restrict(dp dataprovider.DataProvider, perm, collection string) types.Reader {
	return types.ReaderFunc(func(ctx context.Context, userID int, fqfields []string, result map[string]bool) error {
		if len(fqfields) == 0 {
			return nil
		}

		parts := strings.Split(fqfields[0], "/")
		meetingID, err := dp.MeetingFromModel(ctx, collection+"/"+parts[1])
		if err != nil {
			return fmt.Errorf("getting meeting from model: %w", err)
		}

		if err := EnsurePerms(ctx, dp, userID, meetingID, perm); err != nil {
			return nil
		}

		for _, fqfield := range fqfields {
			result[fqfield] = true
		}
		return nil
	})
}
