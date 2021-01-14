package collection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
	"github.com/OpenSlides/openslides-permission-service/internal/perm"
)

// User handels the permissions of user-actions and the user collection.
func User(dp dataprovider.DataProvider) perm.ConnecterFunc {
	u := &user{dp: dp}
	return func(s perm.HandlerStore) {
		s.RegisterWriteHandler("user.create", perm.WriteCheckerFunc(u.create))

		s.RegisterReadHandler("user", perm.ReadCheckerFunc(u.read))
	}
}

type user struct {
	dp dataprovider.DataProvider
}

func (u *user) create(ctx context.Context, userID int, payload map[string]json.RawMessage) (bool, error) {
	var orgaLevel string
	if err := u.dp.GetIfExist(ctx, fmt.Sprintf("user/%d/organisation_management_level", userID), &orgaLevel); err != nil {
		return false, fmt.Errorf("getting organisation level: %w", err)
	}
	switch orgaLevel {
	case "can_manage_organisation", "can_manage_users":
		return true, nil
	default:
		return false, nil
	}
}

func (u *user) read(ctx context.Context, userID int, fqfields []perm.FQField, result map[string]bool) error {
	// Bei can_manage_organisation und can_manage_users(?) kann man alles sehen.

	// TODO: by $-Feldern auch abgeleitete Felder behandeln.

	// Level:
	// 0: no perms
	// 1: can see
	// 2: can see extra
	// 3: can_manage or committtee manager
	meetingLevel := make(map[int]int)

	// Ein committe manager darf alles(?) von usern in seinen meetings sehen.

	// TODO: Ein nutzer darf sich immer selbst sehen, kann man das vorziehen? Nicht wenn er mehr rechte aus einem meeting erhält.

	grouped := groupByID(fqfields)
	for _, fqfields := range grouped {
		var meetingIDs []int
		if err := u.dp.GetIfExist(ctx, fmt.Sprintf("user/%d/is_present_in_meeting_ids", fqfields[0].ID), &meetingIDs); err != nil {
			return fmt.Errorf("getting meeting ids: %w", err)
		}
		for _, meetingID := range meetingIDs {
			level, ok := meetingLevel[meetingID]
			if !ok {
				// TODO: check for committee manager
				perms, err := perm.New(ctx, u.dp, userID, meetingID)
				if err != nil {
					return fmt.Errorf("getting perms for user %d in meeting %d: %w", userID, meetingID, err)
				}
				if perms.Has("user.can_see") {
					level = 1
				}
				if perms.Has("user.can_see_extra_data") {
					level = 2
				}
				if perms.Has("user.can_manage") {
					level = 3
				}

				meetingLevel[meetingID] = level
			}

			// TODO: use level to find out what fields to see
		}
	}
	return nil
}

// groupByID groups a list of fqfields by there id part.
//
// It expects the input to be sorted.
func groupByID(fqfields []perm.FQField) map[int][]perm.FQField {
	grouped := make(map[int][]perm.FQField)
	for _, f := range fqfields {
		grouped[f.ID] = append(grouped[f.ID], f)
	}
	return grouped
}
