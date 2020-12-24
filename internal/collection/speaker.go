package collection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
	"github.com/OpenSlides/openslides-permission-service/internal/perm"
)

// Speaker handels the permissions of the speaker collection.
type Speaker struct {
	dp dataprovider.DataProvider
}

// NewSpeaker initializes a Speaker.
func NewSpeaker(dp dataprovider.DataProvider) *Speaker {
	return &Speaker{
		dp: dp,
	}
}

// Connect connects the list_of_speakers routes.
func (sp *Speaker) Connect(s perm.HandlerStore) {
	s.RegisterWriteHandler("speaker.delete", perm.WriteCheckerFunc(sp.delete))
	s.RegisterReadHandler("speaker", perm.ReadeCheckerFunc(sp.read))

}

func (sp *Speaker) delete(ctx context.Context, userID int, payload map[string]json.RawMessage) (map[string]interface{}, error) {
	fqid := "speaker/" + string(payload["id"])
	var sUserID int
	if err := sp.dp.Get(ctx, fqid+"/user_id", &sUserID); err != nil {
		return nil, fmt.Errorf("getting `%s/user_id` from DB: %w", fqid, err)
	}

	// Speaker is deleting him self.
	if sUserID == userID {
		return nil, nil
	}

	// Check if request-user is list-of-speaker-manager

	meetingID, err := sp.dp.MeetingFromModel(ctx, fqid)
	if err != nil {
		return nil, fmt.Errorf("getting meeting_id from speaker model: %w", err)
	}

	if err := perm.EnsurePerms(ctx, sp.dp, userID, meetingID, "agenda.can_manage_list_of_speakers"); err != nil {
		return nil, fmt.Errorf("ensuring list-of-speaker-manager perms: %w", err)
	}

	return nil, nil
}

func (sp *Speaker) read(ctx context.Context, userID int, fqfields []perm.FQField, result map[string]bool) error {
	var hasPerm bool
	var lastID int
	for _, fqfield := range fqfields {
		if lastID != fqfield.ID {
			hasPerm = sp.hasReadPerm(ctx, userID, fqfield)
		}
		if hasPerm {
			result[fqfield.String()] = true
		}
	}
	return nil
}

func (sp *Speaker) hasReadPerm(ctx context.Context, userID int, fqfield perm.FQField) bool {
	fqid := fmt.Sprintf("speaker/%d", fqfield.ID)
	meetingID, err := sp.dp.MeetingFromModel(ctx, fqid)
	if err != nil {
		return false
	}

	if err := perm.EnsurePerms(ctx, sp.dp, userID, meetingID, "agenda.can_see_list_of_speakers"); err != nil {
		return false
	}
	return true
}