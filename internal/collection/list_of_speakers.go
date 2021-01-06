package collection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
	"github.com/OpenSlides/openslides-permission-service/internal/perm"
)

// ListOfSpeaker connects the handers for list of speakers.
func ListOfSpeaker(dp dataprovider.DataProvider) perm.ConnecterFunc {
	l := &listOfSpeaker{dp}
	return func(s perm.HandlerStore) {
		s.RegisterWriteHandler("list_of_speakers.delete", perm.WriteCheckerFunc(l.delete))
		s.RegisterReadHandler("list_of_speakers", perm.ReadCheckerFunc(l.read))
	}
}

type listOfSpeaker struct {
	dp dataprovider.DataProvider
}

func (l *listOfSpeaker) delete(ctx context.Context, userID int, payload map[string]json.RawMessage) (map[string]interface{}, error) {
	return nil, perm.NotAllowedf("list_of_speaker.delete is an internal action.")
}

func (l *listOfSpeaker) read(ctx context.Context, userID int, fqfields []perm.FQField, result map[string]bool) error {
	return perm.AllFields(fqfields, result, func(fqfield perm.FQField) (bool, error) {
		fqid := fmt.Sprintf("list_of_speakers/%d", fqfield.ID)

		// If the request user is a speaker in the list of speakers, he can see the list.
		var sids []int
		if err := l.dp.Get(ctx, fqid+"/speaker_ids", &sids); err != nil {
			return false, fmt.Errorf("getting speaker object ids: %w", err)
		}

		for _, sid := range sids {
			var suid int
			if err := l.dp.Get(ctx, fmt.Sprintf("speaker/%d/user_id", sid), &suid); err != nil {
				return false, fmt.Errorf("getting speaker user id: %w", err)
			}

			if suid == userID {
				return true, nil
			}
		}

		meetingID, err := l.dp.MeetingFromModel(ctx, fqid)
		if err != nil {
			var errDoesNotExist dataprovider.DoesNotExistError
			if errors.As(err, &errDoesNotExist) {
				return false, nil
			}
			return false, fmt.Errorf("getting meetingID from model %s: %w", fqid, err)
		}

		allowed, err := perm.IsAllowed(perm.EnsurePerm(ctx, l.dp, userID, meetingID, "agenda.can_see_list_of_speakers"))
		if err != nil {
			return false, fmt.Errorf("ensuring perm %w", err)
		}
		return allowed, nil
	})
}
