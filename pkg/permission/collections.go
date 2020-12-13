package permission

import (
	"github.com/OpenSlides/openslides-permission-service/internal/collection"
	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
)

func openSlidesCollections(dp DataProvider) map[string]Collection {
	return map[string]Collection{
		"topic": collection.CreateGeneric(dataprovider.DataProvider{External: dp}, "topic", "agenda.can_manage"),
	}
}
