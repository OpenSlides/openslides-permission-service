package permission

import "github.com/OpenSlides/openslides-permission-service/internal/collection"

func openSlidesCollections(dp DataProvider) map[string]Collection {
	return map[string]Collection{
		"topic": collection.CreateGeneric(dp, "topic", "agenda.can_manage"),
	}
}
