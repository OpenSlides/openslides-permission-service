package permission

import (
	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
	"github.com/OpenSlides/openslides-permission-service/internal/perm"
	"github.com/OpenSlides/openslides-permission-service/internal/perm/agenda"
	"github.com/OpenSlides/openslides-permission-service/internal/perm/autogen"
	"github.com/OpenSlides/openslides-permission-service/internal/perm/user"
)

func openSlidesCollections(edp DataProvider) []perm.Connecter {
	dp := dataprovider.DataProvider{External: edp}

	return []perm.Connecter{
		autogen.NewAutogen(dp),

		agenda.NewSpeaker(dp),
		user.NewPersonalNote(dp),

		//assignment.NewCandidate(dp),
	}
}
