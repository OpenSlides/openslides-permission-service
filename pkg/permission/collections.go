package permission

import (
	"github.com/OpenSlides/openslides-permission-service/internal/collection"
	"github.com/OpenSlides/openslides-permission-service/internal/collection/assignment"
	"github.com/OpenSlides/openslides-permission-service/internal/collection/group"
	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
	"github.com/OpenSlides/openslides-permission-service/internal/types"
)

func openSlidesCollections(edp DataProvider) []types.Connecter {
	dp := dataprovider.DataProvider{External: edp}
	return []types.Connecter{
		collection.NewGeneric(dp, "agenda_item", "agenda.can_see", "agenda.can_manage"),

		collection.NewGeneric(dp, "assignment", "assignments.can_see", "assignments.can_manage"),
		assignment.NewCandidate(dp),

		collection.NewGeneric(dp, "topic", "agenda.can_see", "agenda.can_manage"),

		group.NewGroup(dp),

		// TODO: assignment_poll
		// // TODO: committee

		// "list_of_speakers.update":              listofspeakers.Update,
		// "list_of_speakers.delete_all_speakers": listofspeakers.DeleteAllSpeakers,
		// "list_of_speakers.re_add_last":         listofspeakers.ReAddLast,

		// // TODO: mediafile
		// // TODO: meeting
		// // TODO: motion

		// "motion_block.create": motionblock.Create,
		// "motion_block.update": motionblock.Update,
		// "motion_block.delete": motionblock.Delete,

		// "motion_category.create":                   motioncategory.Create,
		// "motion_category.update":                   motioncategory.Update,
		// "motion_category.delete":                   motioncategory.Delete,
		// "motion_category.sort":                     motioncategory.Sort,
		// "motion_category.sort_motions_in_category": motioncategory.SortMotionsInCategory,
		// "motion_category.number_motions":           motioncategory.NumberMotions,

		// "motion_change_recommendation.create": motion_change_recommendation.Create,
		// "motion_change_recommendation.update": motion_change_recommendation.Update,
		// "motion_change_recommendation.delete": motion_change_recommendation.Delete,

		// // TODO: motion_comment

		// "motion_comment_section.create": motioncommentsection.Create,
		// "motion_comment_section.update": motioncommentsection.Update,
		// "motion_comment_section.delete": motioncommentsection.Delete,
		// // TODO: sort

		// // TODO: motion_poll
		// // TODO: motion_state

		// "motion_statute_paragraph.create": motion_statute_paragraph.Create,
		// "motion_statute_paragraph.update": motion_statute_paragraph.Update,
		// "motion_statute_paragraph.delete": motion_statute_paragraph.Delete,
		// // TODO: sort

		// // TODO: motion_submitter

		// "motion_workflow.create": motionworkflow.Create,
		// "motion_workflow.update": motionworkflow.Update,
		// "motion_workflow.delete": motionworkflow.Delete,

		// // TODO: personal_note
		// // TODO: speaker

		// "tag.create": tag.Create,
		// "tag.update": tag.Update,
		// "tag.delete": tag.Delete,

		// // TODO: users
	}
}
