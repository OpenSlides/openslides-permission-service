package collection

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
	"github.com/OpenSlides/openslides-permission-service/internal/perm"
)

// Motion handels permissions of motions objects.
type Motion struct {
	dp dataprovider.DataProvider
}

// NewMotion initializes a motion.
func NewMotion(dp dataprovider.DataProvider) *Motion {
	return &Motion{
		dp: dp,
	}
}

// Connect registers the Motion handlers.
func (m *Motion) Connect(s perm.HandlerStore) {
	s.RegisterWriteHandler("motion.delete", m.modify("motion.can_manage"))
	s.RegisterWriteHandler("motion.set_state", m.modify("motion.can_manage_metadata"))
	s.RegisterWriteHandler("motion.create", m.create())

	s.RegisterReadHandler("motion", perm.ReadeCheckerFunc(m.readMotion))
	s.RegisterReadHandler("motion_submitter", perm.ReadeCheckerFunc(m.readMotionSubmitter))
}

func (m *Motion) create() perm.WriteCheckerFunc {
	allowList := map[string]bool{
		"title":                true,
		"text":                 true,
		"reason":               true,
		"category_id":          true,
		"statute_paragraph_id": true,
		"workflow_id":          true,
		"meeting_id":           true,
	}

	allowListAmendment := map[string]bool{
		"parent_id":            true,
		"amendment_paragraphs": true,
	}

	for k := range allowList {
		allowListAmendment[k] = true
	}

	return func(ctx context.Context, userID int, payload map[string]json.RawMessage) (map[string]interface{}, error) {
		meetingID, _ := strconv.Atoi(string(payload["meeting_id"]))

		perms, err := perm.Perms(ctx, userID, meetingID, m.dp)
		if err != nil {
			return nil, fmt.Errorf("fetching perms: %w", err)
		}

		if perms.HasOne("motion.can_manage") {
			return nil, nil
		}

		requiredPerm := "motion.can_create"
		aList := allowList
		if _, ok := payload["parent_id"]; ok {
			requiredPerm = "motion.can_create_amendment"
			aList = allowListAmendment
		}

		if !perms.HasOne(requiredPerm) {
			return nil, perm.NotAllowedf("User %d does not have permission %s", userID, requiredPerm)
		}

		for e := range payload {
			if !aList[string(e)] {
				return nil, perm.NotAllowedf("Field `%s` is forbidden for non manager.", e)
			}
		}

		return nil, nil
	}
}

func (m *Motion) modify(managePerm string) perm.WriteCheckerFunc {
	return func(ctx context.Context, userID int, payload map[string]json.RawMessage) (map[string]interface{}, error) {
		motionFQID := fmt.Sprintf("motion/%s", payload["id"])
		meetingID, err := m.dp.MeetingFromModel(ctx, motionFQID)
		if err != nil {
			return nil, fmt.Errorf("getting meeting for %s: %w", motionFQID, err)
		}

		isManager, err := perm.IsAllowed(perm.EnsurePerm(ctx, m.dp, userID, meetingID, managePerm))
		if err != nil {
			return nil, fmt.Errorf("checking meta manager permission: %w", err)
		}

		if isManager {
			return nil, nil
		}

		var submitterIDs []int
		if err := m.dp.Get(ctx, motionFQID+"/submitter_ids", &submitterIDs); err != nil {
			return nil, fmt.Errorf("getting submitter ids: %w", err)
		}

		var isSubmitter bool
		for _, sid := range submitterIDs {
			var sUserID int
			if err := m.dp.Get(ctx, fmt.Sprintf("motion_submitter/%d/user_id", sid), &sUserID); err != nil {
				return nil, fmt.Errorf("getting userid of sumitter %d: %w", sid, err)
			}
			if sUserID == userID {
				isSubmitter = true
				break
			}
		}

		if !isSubmitter {
			return nil, perm.NotAllowedf("User %d is not a manager and not a submitter of %s", userID, motionFQID)
		}

		var stateID int
		if err := m.dp.Get(ctx, motionFQID+"/state_id", &stateID); err != nil {
			return nil, fmt.Errorf("getting stateID: %w", err)
		}

		var allowSubmitterEdit bool
		if err := m.dp.Get(ctx, fmt.Sprintf("motion_state/%d/allow_submitter_edit", stateID), &allowSubmitterEdit); err != nil {
			return nil, fmt.Errorf("getting allow_submitter_edit: %w", err)
		}

		if !allowSubmitterEdit {
			return nil, perm.NotAllowedf("Motion state does not allow submitter edites")
		}

		return nil, nil
	}
}

func (m *Motion) canSeeMotion(ctx context.Context, userID int, motionID int) (bool, error) {
	motionFQID := fmt.Sprintf("motion/%d", motionID)
	meetingID, err := m.dp.MeetingFromModel(ctx, motionFQID)
	if err != nil {
		return false, fmt.Errorf("getting meetingID from model %s: %w", motionFQID, err)
	}

	committeeID, err := m.dp.CommitteeID(ctx, meetingID)
	if err != nil {
		return false, fmt.Errorf("getting committee id for meeting: %w", err)
	}

	committeeManager, err := m.dp.IsManager(ctx, userID, committeeID)
	if err != nil {
		return false, fmt.Errorf("check for manager: %w", err)
	}
	if committeeManager {
		return true, nil
	}

	isMeeting, err := m.dp.InMeeting(ctx, userID, meetingID)
	if err != nil {
		return false, fmt.Errorf("Looking for user %d in meeting %d: %w", userID, meetingID, err)
	}
	if !isMeeting {
		return false, nil
	}

	perms, err := perm.Perms(ctx, userID, meetingID, m.dp)
	if err != nil {
		return false, fmt.Errorf("getting user permissions: %w", err)
	}

	if perms.HasOne("motion.can_manage") {
		return true, nil
	}

	if !perms.HasOne("motion.can_see") {
		return false, nil
	}

	var stateID int
	if err := m.dp.Get(ctx, motionFQID+"/state_id", &stateID); err != nil {
		return false, fmt.Errorf("getting field %s/state_id: %w", motionFQID, err)
	}

	var restriction []string
	field := fmt.Sprintf("motion_state/%d/restrictions", stateID)
	if err := m.dp.Get(ctx, field, &restriction); err != nil {
		return false, fmt.Errorf("getting field %s: %w", field, err)
	}

	if len(restriction) == 0 {
		return true, nil
	}

	for _, r := range restriction {
		switch r {
		case "motion.can_see_internal", "motion.can_manage_metadata", "motion.can_manage":
			if perms.HasOne(r) {
				return true, nil
			}

		case "is_submitter":
			var submitterIDs []int
			if err := m.dp.Get(ctx, motionFQID+"/submitter_ids", &submitterIDs); err != nil {
				return false, fmt.Errorf("getting field %s/submitter_ids: %w", motionFQID, err)
			}

			for _, sid := range submitterIDs {
				var uid int
				f := fmt.Sprintf("motion_submitter/%d/user_id", sid)
				if err := m.dp.Get(ctx, f, &uid); err != nil {
					return false, fmt.Errorf("getting field %s: %w", f, err)
				}
				if uid == userID {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (m *Motion) readMotion(ctx context.Context, userID int, fqfields []perm.FQField, result map[string]bool) error {
	return perm.AllFields(fqfields, result, func(fqfield perm.FQField) (bool, error) {
		return m.canSeeMotion(ctx, userID, fqfield.ID)
	})
}

func (m *Motion) readMotionSubmitter(ctx context.Context, userID int, fqfields []perm.FQField, result map[string]bool) error {
	return perm.AllFields(fqfields, result, func(fqfield perm.FQField) (bool, error) {
		var motionID int
		if err := m.dp.Get(ctx, fmt.Sprintf("motion_submitter/%d/motion_id", fqfield.ID), &motionID); err != nil {
			return false, fmt.Errorf("getting motionID: %w", err)
		}
		return m.canSeeMotion(ctx, userID, motionID)
	})
}
