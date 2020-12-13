package permission

import (
	"github.com/OpenSlides/openslides-permission-service/internal/collection"
	"github.com/OpenSlides/openslides-permission-service/internal/collection/assignment"
	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
)

func openSlidesCollections(edp DataProvider) map[string]Collection {
	dp := dataprovider.DataProvider{External: edp}
	return map[string]Collection{
		"agenda_item":          collection.CreateGeneric(dp, "agenda_item", "agenda.can_manage"), // TODO: assign, sort
		"assignment":           collection.CreateGeneric(dp, "assignment", "assignments.can_manage"),
		"assignment_candidate": assignment.NewCandidate(dp),

		"topic": collection.CreateGeneric(dp, "topic", "agenda.can_manage"),
	}
}
