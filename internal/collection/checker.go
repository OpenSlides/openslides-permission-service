package collection

import (
	"context"
	"fmt"

	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
)

// EnsurePerms makes sure the user has at least one of the given permissions.
func EnsurePerms(ctx context.Context, dp dataprovider.DataProvider, userID int, meetingID int, permissions ...string) error {
	committeeID, err := dp.CommitteeID(ctx, meetingID)
	if err != nil {
		return fmt.Errorf("getting committee id for meeting: %w", err)
	}

	committeeManager, err := dp.IsManager(ctx, userID, committeeID)
	if err != nil {
		return fmt.Errorf("check for manager: %w", err)
	}
	if committeeManager {
		return nil
	}

	canSeeMeeting, err := dp.InMeeting(ctx, userID, meetingID)
	if err != nil {
		return err
	}
	if !canSeeMeeting {
		return NotAllowedf("User %d is not in meeting %d", userID, meetingID)
	}

	perms, err := Perms(ctx, userID, meetingID, dp)
	if err != nil {
		return fmt.Errorf("getting user permissions: %w", err)
	}

	hasPerms := perms.HasOne(permissions...)
	if !hasPerms {
		return NotAllowedf("User %d has not the required permission in meeting %d", userID, meetingID)
	}

	return nil
}
