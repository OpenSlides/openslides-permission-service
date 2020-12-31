package permission

import (
	"github.com/OpenSlides/openslides-permission-service/internal/collection"
	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
	"github.com/OpenSlides/openslides-permission-service/internal/perm"
)

func openSlidesCollections(edp DataProvider) []perm.Connecter {
	dp := dataprovider.DataProvider{External: edp}

	return []perm.Connecter{
		collection.NewAutogen(dp),

		collection.NewAgendaItem(dp),
		collection.NewSpeaker(dp),
		collection.NewPersonalNote(dp),
		collection.NewGroup(dp),
		collection.ReadPerm(dp, "list_of_speakers", "agenda.can_see_list_of_speakers"),
		collection.ReadPerm(dp, "assignment", "assingment.can_see"),
		collection.ReadPerm(dp, "assignment_candidate", "assingment.can_see"),
		collection.ReadInMeeting(dp, "tag"),
		collection.ReadPerm(dp, "projector_message", "meeting.can_see_projector"),
		collection.ReadPerm(dp, "projector_countdown", "meeting.can_see_projector"),
		collection.NewMotion(dp),
	}
}
