package assignment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/OpenSlides/openslides-permission-service/internal/collection"
	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
)

// Candidate is the collection for assignment candidates.
type Candidate struct {
	collection.Generic
}

// NewCandidate creates a new AssignmentCandidate collection.
func NewCandidate(dp dataprovider.DataProvider) *Candidate {
	return new(Candidate)
}

// IsAllowed tells, if the user has the perm to modify the candidate.
func (c *Candidate) IsAllowed(ctx context.Context, name string, userID int, data map[string]json.RawMessage) (map[string]interface{}, error) {
	return nil, fmt.Errorf("TODO")
}
