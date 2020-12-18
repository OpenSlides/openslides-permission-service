package permission

import (
	"github.com/OpenSlides/openslides-permission-service/internal/collection"
	"github.com/OpenSlides/openslides-permission-service/internal/collection/assignment"
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
	}
}
