package user

import (
	"github.com/OpenSlides/openslides-permission-service/internal/collection"
	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
	"github.com/OpenSlides/openslides-permission-service/internal/types"
)

// Group is the collection for user groups.
type Group struct {
	dp dataprovider.DataProvider
}

// NewGroup creates a new Group collection.
func NewGroup(dp dataprovider.DataProvider) *Group {
	return &Group{
		dp: dp,
	}
}

// Connect connects the assignment_candidate routes.
func (g *Group) Connect(s types.HandlerStore) {
	gen := collection.NewGeneric(g.dp, "group", "users.can_see", "users.can_manage")
	gen.Connect(s)

	s.RegisterWriteHandler("group.set_permission", collection.Modify(g.dp, "users.can_manate", "group"))
}
