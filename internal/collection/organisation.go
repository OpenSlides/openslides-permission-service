package collection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/OpenSlides/openslides-permission-service/internal/dataprovider"
	"github.com/OpenSlides/openslides-permission-service/internal/perm"
)

// Organisation handels permissions for the organisation.
func Organisation(dp dataprovider.DataProvider) perm.ConnecterFunc {
	return func(s perm.HandlerStore) {
		s.RegisterAction("organisation.update", perm.ActionCheckerFunc(
			func(ctx context.Context, userID int, payload map[string]json.RawMessage) (bool, error) {
				var orgaLevel string
				if err := dp.GetIfExist(ctx, fmt.Sprintf("user/%d/organisation_management_level", userID), &orgaLevel); err != nil {
					return false, fmt.Errorf("getting organisation level: %w", err)
				}

				if orgaLevel == "can_manage_organisation" {
					return true, nil
				}
				return false, nil
			},
		))
	}
}
