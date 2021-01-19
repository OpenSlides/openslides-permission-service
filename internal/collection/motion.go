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

	s.RegisterReadHandler("motion", perm.ReadCheckerFunc(m.readMotion))
	s.RegisterReadHandler("motion_submitter", perm.ReadCheckerFunc(m.readSubmitter))
	s.RegisterReadHandler("motion_block", m.readBlock())
	s.RegisterReadHandler("motion_change_recommendation", m.readChangeRecommendation())
	s.RegisterReadHandler("motion_comment_section", perm.ReadCheckerFunc(m.readCommentSection))
	s.RegisterReadHandler("motion_comment", perm.ReadCheckerFunc(m.readComment))
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

	return func(ctx context.Context, userID int, payload map[string]json.RawMessage) (bool, error) {
		meetingID, err := strconv.Atoi(string(payload["meeting_id"]))
		if err != nil {
			return false, fmt.Errorf("invalid field meeting_id in payload: %w", err)
		}

		perms, err := perm.New(ctx, m.dp, userID, meetingID)
		if err != nil {
			return false, fmt.Errorf("fetching perms: %w", err)
		}

		if perms.Has("motion.can_manage") {
			return true, nil
		}

		requiredPerm := "motion.can_create"
		aList := allowList
		if _, ok := payload["parent_id"]; ok {
			requiredPerm = "motion.can_create_amendment"
			aList = allowListAmendment
		}

		if !perms.Has(requiredPerm) {
			perm.LogNotAllowedf("User %d does not have permission %s", userID, requiredPerm)
			return false, nil
		}

		for e := range payload {
			if !aList[string(e)] {
				perm.LogNotAllowedf("Field `%s` is forbidden for non manager.", e)
				return false, nil
			}
		}

		return true, nil
	}
}

func (m *Motion) modify(managePerm string) perm.WriteCheckerFunc {
	return func(ctx context.Context, userID int, payload map[string]json.RawMessage) (bool, error) {
		motionFQID := fmt.Sprintf("motion/%s", payload["id"])
		meetingID, err := m.dp.MeetingFromModel(ctx, motionFQID)
		if err != nil {
			return false, fmt.Errorf("getting meeting for %s: %w", motionFQID, err)
		}

		isManager, err := perm.HasPerm(ctx, m.dp, userID, meetingID, managePerm)
		if err != nil {
			return false, fmt.Errorf("checking meta manager permission: %w", err)
		}

		if isManager {
			return true, nil
		}

		var submitterIDs []int
		if err := m.dp.GetIfExist(ctx, motionFQID+"/submitter_ids", &submitterIDs); err != nil {
			return false, fmt.Errorf("getting submitter ids: %w", err)
		}

		var isSubmitter bool
		for _, sid := range submitterIDs {
			var sUserID int
			if err := m.dp.Get(ctx, fmt.Sprintf("motion_submitter/%d/user_id", sid), &sUserID); err != nil {
				return false, fmt.Errorf("getting userid of sumitter %d: %w", sid, err)
			}
			if sUserID == userID {
				isSubmitter = true
				break
			}
		}

		if !isSubmitter {
			perm.LogNotAllowedf("User %d is not a manager and not a submitter of %s", userID, motionFQID)
			return false, nil
		}

		var stateID int
		if err := m.dp.Get(ctx, motionFQID+"/state_id", &stateID); err != nil {
			return false, fmt.Errorf("getting stateID: %w", err)
		}

		var allowSubmitterEdit bool
		if err := m.dp.GetIfExist(ctx, fmt.Sprintf("motion_state/%d/allow_submitter_edit", stateID), &allowSubmitterEdit); err != nil {
			return false, fmt.Errorf("getting allow_submitter_edit: %w", err)
		}

		if !allowSubmitterEdit {
			perm.LogNotAllowedf("Motion state does not allow submitter edites")
			return false, nil
		}

		return true, nil
	}
}

func canSeeMotion(ctx context.Context, dp dataprovider.DataProvider, userID int, motionID int, perms *perm.Permission) (bool, error) {
	if perms.Has("motion.can_manage") {
		return true, nil
	}

	if !perms.Has("motion.can_see") {
		return false, nil
	}

	motionFQID := fmt.Sprintf("motion/%d", motionID)

	var stateID int
	if err := dp.Get(ctx, motionFQID+"/state_id", &stateID); err != nil {
		return false, fmt.Errorf("getting motion state: %w", err)
	}

	var restriction []string
	field := fmt.Sprintf("motion_state/%d/restrictions", stateID)
	if err := dp.GetIfExist(ctx, field, &restriction); err != nil {
		return false, fmt.Errorf("getting field %s: %w", field, err)
	}

	if len(restriction) == 0 {
		return true, nil
	}

	for _, r := range restriction {
		switch r {
		case "motion.can_see_internal", "motion.can_manage_metadata", "motion.can_manage":
			if perms.Has(r) {
				return true, nil
			}

		case "is_submitter":
			var submitterIDs []int
			if err := dp.GetIfExist(ctx, motionFQID+"/submitter_ids", &submitterIDs); err != nil {
				return false, fmt.Errorf("getting field %s/submitter_ids: %w", motionFQID, err)
			}

			for _, sid := range submitterIDs {
				var uid int
				f := fmt.Sprintf("motion_submitter/%d/user_id", sid)
				if err := dp.Get(ctx, f, &uid); err != nil {
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
		meetingID, err := m.dp.MeetingFromModel(ctx, fmt.Sprintf("motion/%d", fqfield.ID))
		if err != nil {
			return false, fmt.Errorf("getting meetingID from motion: %w", err)
		}

		perms, err := perm.New(ctx, m.dp, userID, meetingID)
		if err != nil {
			return false, fmt.Errorf("getting user permissions: %w", err)
		}

		return canSeeMotion(ctx, m.dp, userID, fqfield.ID, perms)
	})
}

func (m *Motion) readSubmitter(ctx context.Context, userID int, fqfields []perm.FQField, result map[string]bool) error {
	return perm.AllFields(fqfields, result, func(fqfield perm.FQField) (bool, error) {
		var motionID int
		if err := m.dp.Get(ctx, fmt.Sprintf("motion_submitter/%d/motion_id", fqfield.ID), &motionID); err != nil {
			return false, fmt.Errorf("getting motionID: %w", err)
		}

		meetingID, err := m.dp.MeetingFromModel(ctx, fmt.Sprintf("motion/%d", motionID))
		if err != nil {
			return false, fmt.Errorf("getting meetingID from motion: %w", err)
		}

		perms, err := perm.New(ctx, m.dp, userID, meetingID)
		if err != nil {
			return false, fmt.Errorf("getting user permissions: %w", err)
		}
		return canSeeMotion(ctx, m.dp, userID, motionID, perms)
	})
}

func (m *Motion) readBlock() perm.ReadCheckerFunc {
	return func(ctx context.Context, userID int, fqfields []perm.FQField, result map[string]bool) error {
		return perm.AllFields(fqfields, result, func(fqfield perm.FQField) (bool, error) {
			fqid := fmt.Sprintf("motion_block/%d", fqfield.ID)
			meetingID, err := m.dp.MeetingFromModel(ctx, fqid)
			if err != nil {
				return false, fmt.Errorf("getting meetingID from model %s: %w", fqid, err)
			}

			perms, err := perm.New(ctx, m.dp, userID, meetingID)
			if err != nil {
				return false, fmt.Errorf("getting user permissions: %w", err)
			}

			if perms.Has("motion.can_manage") {
				return true, nil
			}

			if !perms.Has("motion.can_see") {
				return false, nil
			}

			var internal bool
			if err := m.dp.GetIfExist(ctx, fqid+"/internal", &internal); err != nil {
				return false, fmt.Errorf("get /internal: %w", err)
			}

			if !internal {
				return true, nil
			}

			return false, nil
		})
	}
}

func (m *Motion) readChangeRecommendation() perm.ReadCheckerFunc {
	return func(ctx context.Context, userID int, fqfields []perm.FQField, result map[string]bool) error {
		return perm.AllFields(fqfields, result, func(fqfield perm.FQField) (bool, error) {
			fqid := fmt.Sprintf("motion_change_recommendation/%d", fqfield.ID)
			meetingID, err := m.dp.MeetingFromModel(ctx, fqid)
			if err != nil {
				return false, fmt.Errorf("getting meetingID from model %s: %w", fqid, err)
			}

			perms, err := perm.New(ctx, m.dp, userID, meetingID)
			if err != nil {
				return false, fmt.Errorf("getting user permissions: %w", err)
			}

			if perms.Has("motion.can_manage") {
				return true, nil
			}

			var motionID int
			if err := m.dp.Get(ctx, fqid+"/motion_id", &motionID); err != nil {
				return false, fmt.Errorf("getting motion id: %w", err)
			}

			motionOK, err := canSeeMotion(ctx, m.dp, userID, motionID, perms)
			if err != nil {
				return false, fmt.Errorf("checking permission for motion: %w", err)
			}
			if !motionOK {
				return false, nil
			}

			var internal bool
			if err := m.dp.GetIfExist(ctx, fqid+"/internal", &internal); err != nil {
				return false, fmt.Errorf("get /internal: %w", err)
			}

			if !internal {
				return true, nil
			}

			return perms.Has("motion.can_manage"), nil
		})
	}
}

func (m *Motion) canSeeCommentSection(ctx context.Context, userID, id int) (bool, error) {
	fqid := fmt.Sprintf("motion_comment_section/%d", id)
	meetingID, err := m.dp.MeetingFromModel(ctx, fqid)
	if err != nil {
		return false, fmt.Errorf("getting meetingID from model %s: %w", fqid, err)
	}

	perms, err := perm.New(ctx, m.dp, userID, meetingID)
	if err != nil {
		return false, fmt.Errorf("getting user permissions: %w", err)
	}

	if perms.Has("motion.can_manage") {
		return true, nil
	}

	var motionID int
	if err := m.dp.Get(ctx, fqid+"/motion_id", &motionID); err != nil {
		return false, fmt.Errorf("getting motion id: %w", err)
	}

	motionOK, err := canSeeMotion(ctx, m.dp, userID, motionID, perms)
	if err != nil {
		return false, fmt.Errorf("checking permission for motion: %w", err)
	}
	if !motionOK {
		return false, nil
	}

	var readGroupIDs []int
	if err := m.dp.GetIfExist(ctx, fqid+"/read_group_ids", &readGroupIDs); err != nil {
		return false, fmt.Errorf("getting read groups: %w", err)
	}
	for _, uid := range readGroupIDs {
		if uid == userID {
			return true, nil
		}
	}
	return false, nil
}

func (m *Motion) readCommentSection(ctx context.Context, userID int, fqfields []perm.FQField, result map[string]bool) error {
	return perm.AllFields(fqfields, result, func(fqfield perm.FQField) (bool, error) {
		return m.canSeeCommentSection(ctx, userID, fqfield.ID)
	})
}

func (m *Motion) readComment(ctx context.Context, userID int, fqfields []perm.FQField, result map[string]bool) error {
	return perm.AllFields(fqfields, result, func(fqfield perm.FQField) (bool, error) {
		var sectionID int
		if err := m.dp.Get(ctx, fmt.Sprintf("motion_comment/%d/section_id", fqfield.ID), &sectionID); err != nil {
			return false, fmt.Errorf("getting section id: %w", err)
		}
		return m.canSeeCommentSection(ctx, userID, sectionID)
	})
}

func canSeeMotionSupporter(ctx context.Context, dp dataprovider.DataProvider, userID int, p *perm.Permission, ids []int) (bool, error) {
	for _, id := range ids {
		b, err := canSeeMotion(ctx, dp, userID, id, p)
		if err != nil {
			return false, err
		}
		if b {
			return true, nil
		}
	}
	return false, nil
}
