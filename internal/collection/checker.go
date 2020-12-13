package collection

import (
	"fmt"
	"strconv"

	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
)

// IsSuperuser returns true, if the user is a superuser. If the user does not
// exist at all, an error is returned.
func IsSuperuser(userID int, dp dataprovider.DataProvider) (bool, error) {
	exists, err := DoesUserExists(userID, dp)
	if err != nil {
		return false, fmt.Errorf("check if user exist: %w", err)
	}
	if !exists {
		return false, NotAllowedf("The user with id %d does not exist", userID)
	}

	superadmin, err := HasUserSuperadminRole(userID, dp)
	if err != nil {
		return false, fmt.Errorf("check for super user role: %w", err)
	}
	return superadmin, nil
}

// DoesUserExists returns true, if an user exist. Returns allways true for
// userID 0.
func DoesUserExists(userID int, dp dataprovider.DataProvider) (bool, error) {
	if userID == 0 {
		return true, nil
	}

	exists, err := DoesModelExists("user/"+strconv.Itoa(userID), dp)
	if err != nil {
		return false, fmt.Errorf("lockup user: %w", err)
	}
	return exists, nil
}

// DoesModelExists returns true, if an object exists in the datastore.
func DoesModelExists(fqid string, dp dataprovider.DataProvider) (bool, error) {
	exists, err := dp.Exists(fqid + "/" + "id")
	if err != nil {
		return false, fmt.Errorf("checking for model existing: %w", err)
	}
	return exists, nil
}

// HasUserSuperadminRole returns true, if the user is in the superuser group.
func HasUserSuperadminRole(userID int, dp dataprovider.DataProvider) (bool, error) {
	// The anonymous is never a superadmin.
	if userID == 0 {
		return false, nil
	}

	// Get superadmin role id.
	superadminRoleID, err := dp.GetInt("organisation/1/superadmin_role_id")
	if err != nil {
		return false, fmt.Errorf("getting superadmin role id: %w", err)
	}

	// Get users role id.
	fqfield := "user/" + strconv.Itoa(userID) + "/role_id"
	if exists, err := dp.Exists(fqfield); !exists || err != nil {
		// The user has no role.
		return false, err
	}

	userRoleID, err := dp.GetInt(fqfield)
	if err != nil {
		return false, fmt.Errorf("getting role_id: %w", err)
	}

	return superadminRoleID == userRoleID, nil
}

// EnsurePerms makes sure the user has at least one of the given permissions.
func EnsurePerms(dp dataprovider.DataProvider, userID int, meetingID int, permissions ...string) error {
	committeeID, err := CommitteeID(meetingID, dp)
	if err != nil {
		return fmt.Errorf("getting committee id for meeting: %w", err)
	}

	committeeManager, err := IsManager(userID, committeeID, dp)
	if err != nil {
		return fmt.Errorf("check for manager: %w", err)
	}
	if committeeManager {
		return nil
	}

	canSeeMeeting, err := InMeeting(userID, meetingID, dp)
	if err != nil {
		return err
	}
	if !canSeeMeeting {
		return NotAllowedf("User %d is not in meeting %d", userID, meetingID)
	}

	perms, err := Perms(userID, meetingID, dp)
	if err != nil {
		return fmt.Errorf("getting user permissions: %w", err)
	}

	hasPerms := perms.HasOne(permissions...)
	if !hasPerms {
		return NotAllowedf("User %d has not the required permission in meeting %d", userID, meetingID)
	}

	return nil
}

// CommitteeID returns the id of a committee from an meeting id.
func CommitteeID(meetingID int, dp dataprovider.DataProvider) (int, error) {
	committeeID, err := dp.GetInt("meeting/" + strconv.Itoa(meetingID) + "/committee_id")
	if err != nil {
		return 0, fmt.Errorf("getting committee id: %w", err)
	}
	return committeeID, nil
}

// IsManager returns true, if the user is a manager in the committee.
func IsManager(userID, committeeID int, dp dataprovider.DataProvider) (bool, error) {
	// The anonymous is never a manager.
	if userID == 0 {
		return false, nil
	}

	// Get committee manager_ids.
	managerIDs, err := dp.GetIntArrayWithDefault("committee/"+strconv.Itoa(committeeID)+"/manager_ids", []int{})
	if err != nil {
		return false, fmt.Errorf("getting committee ids: %w", err)
	}

	for _, id := range managerIDs {
		if userID == id {
			return true, nil
		}
	}
	return false, nil
}

// InMeeting returns true, if the user is in the user_ids list or anonymous.
func InMeeting(userID, meetingID int, dp dataprovider.DataProvider) (bool, error) {
	if userID == 0 {
		enableAnonymous, err := dp.GetBoolWithDefault("meeting/"+strconv.Itoa(meetingID)+"/enable_anonymous", false)
		if err != nil {
			return false, fmt.Errorf("checking anonymous enabled: %w", err)
		}
		return enableAnonymous, nil
	}

	userIds, err := dp.GetIntArrayWithDefault("meeting/"+strconv.Itoa(meetingID)+"/user_ids", []int{})
	if err != nil {
		return false, fmt.Errorf("getting meeting/user_ids: %w", err)
	}

	for _, id := range userIds {
		if id == userID {
			return true, nil
		}
	}
	return false, nil
}

// MeetingFromModel returns the meeting id for an model.
func MeetingFromModel(fqid string, dp dataprovider.DataProvider) (int, error) {
	id, err := dp.GetInt(fqid + "/meeting_id")
	if err != nil {
		return 0, fmt.Errorf("getting meeting id: %w", err)
	}
	return id, nil
}
