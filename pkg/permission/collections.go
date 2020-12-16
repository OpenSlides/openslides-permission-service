package permission

import (
	"github.com/OpenSlides/openslides-permission-service/internal/collection"
	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
)

func openSlidesCollections(edp DataProvider) []Collection {
	dp := dataprovider.DataProvider{External: edp}
	return []Collection{
		collection.NewGeneric(dp, "topic", "agenda.can_see", "agenda.can_manage"),
		collection.NewGeneric(dp, "agenda_item", "agenda.can_see", "agenda.can_manage"),
		collection.NewGeneric(dp, "assignment", "assignments.can_see", "assignments.can_manage"),
	}
}
