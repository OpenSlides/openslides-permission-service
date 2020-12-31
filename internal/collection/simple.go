package collection

import (
	"context"
	"fmt"

	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
	"github.com/OpenSlides/openslides-permission-service/internal/perm"
)

// ReadPerm checks that the user has one permission in a meeting.
//
// If the user has the permission, all fields of the given collection can be
// seen.
func ReadPerm(dp dataprovider.DataProvider, collection string, permission string) perm.ConnecterFunc {
	read := func(ctx context.Context, userID int, fqfields []perm.FQField, result map[string]bool) error {
		return perm.AllFields(fqfields, result, func(fqfield perm.FQField) (bool, error) {
			fqid := fmt.Sprintf("%s/%d", collection, fqfield.ID)
			meetingID, err := dp.MeetingFromModel(ctx, fqid)
			if err != nil {
				return false, fmt.Errorf("getting meetingID from model %s: %w", fqid, err)
			}

			allowed, err := perm.IsAllowed(perm.EnsurePerm(ctx, dp, userID, meetingID, permission))
			if err != nil {
				return false, fmt.Errorf("ensuring perm %w", err)
			}
			return allowed, nil
		})
	}

	return func(s perm.HandlerStore) {
		s.RegisterReadHandler(collection, perm.ReadeCheckerFunc(read))
	}
}

// ReadInMeeting lets the user see the collection, if he is in the meeting.
func ReadInMeeting(dp dataprovider.DataProvider, collection string) perm.ConnecterFunc {
	read := func(ctx context.Context, userID int, fqfields []perm.FQField, result map[string]bool) error {
		return perm.AllFields(fqfields, result, func(fqfield perm.FQField) (bool, error) {
			fqid := fmt.Sprintf("%s/%d", collection, fqfield.ID)
			meetingID, err := dp.MeetingFromModel(ctx, fqid)
			if err != nil {
				return false, fmt.Errorf("getting meetingID from model %s: %w", fqid, err)
			}

			inMeeting, err := dp.InMeeting(ctx, userID, meetingID)
			if err != nil {
				return false, fmt.Errorf("see if user %d is in meeting %d: %w", userID, meetingID, err)
			}
			return inMeeting, nil
		})
	}

	return func(s perm.HandlerStore) {
		s.RegisterReadHandler(collection, perm.ReadeCheckerFunc(read))
	}
}
