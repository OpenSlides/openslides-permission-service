package collection

import (
	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
	"github.com/OpenSlides/openslides-permission-service/internal/types"
)

// Generic is a helper object to create a collection with usual functions.
type Generic struct {
	dp         dataprovider.DataProvider
	collection string
	managePerm string
	readPerm   string
}

// NewGeneric creates a generic collection.
func NewGeneric(dp dataprovider.DataProvider, collection string, readPerm, managePerm string) *Generic {
	return &Generic{
		dp:         dp,
		collection: collection,
		managePerm: managePerm,
		readPerm:   readPerm,
	}
}

// Connect sets the generic routs to the given reader and writer.
func (g *Generic) Connect(s types.HandlerStore) {
	s.RegisterWriteHandler(g.collection+".create", Create(g.dp, g.managePerm, g.collection))
	s.RegisterWriteHandler(g.collection+".update", Modify(g.dp, g.managePerm, g.collection))
	s.RegisterWriteHandler(g.collection+".delete", Modify(g.dp, g.managePerm, g.collection))

	s.RegisterReadHandler(g.collection, Restrict(g.dp, g.readPerm, g.collection))
}
