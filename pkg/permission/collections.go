package permission

import (
	"github.com/OpenSlides/openslides-permission-service/internal/collection"
	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
	"github.com/OpenSlides/openslides-permission-service/internal/perm"
	"github.com/OpenSlides/openslides-permission-service/internal/perm/autogen"
)

func openSlidesCollections(edp DataProvider) []perm.Connecter {
	dp := dataprovider.DataProvider{External: edp}

	return []perm.Connecter{
		autogen.NewAutogen(dp),

		collection.NewSpeaker(dp),
		collection.NewPersonalNote(dp),

		//assignment.NewCandidate(dp),
	}
}
