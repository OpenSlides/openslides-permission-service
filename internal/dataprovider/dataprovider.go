package dataprovider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/OpenSlides/openslides-permission-service/internal/definitions"
)

type externalDataProvider interface {
	// If a field does not exist, it is not returned.
	Get(ctx context.Context, fields ...definitions.Fqfield) ([]json.RawMessage, error)
}

// DataProvider is a wrapper around permission.DataProvider that provides some
// helper functions.
type DataProvider struct {
	External externalDataProvider
}

func (dp *DataProvider) externalGet(ctx context.Context, fields ...definitions.Fqfield) ([]json.RawMessage, error) {
	return dp.External.Get(ctx, fields...)
}

// Get returns a value from the datastore and unpacks it in to the argument value.
//
// The argument value has to be an non nil pointer.
func (dp *DataProvider) Get(ctx context.Context, fqfield string, value interface{}) error {
	fields, err := dp.externalGet(ctx, fqfield)
	if err != nil {
		return fmt.Errorf("getting data from datastore: %w", err)
	}

	if fields[0] == nil {
		return doesNotExistError(fqfield)
	}

	if err := json.Unmarshal(fields[0], value); err != nil {
		return fmt.Errorf("unpacking value: %w", err)
	}
	return nil
}

// GetIfExist behaves like Get() but does not throw an error if the fqfield does
// not exist.
func (dp DataProvider) GetIfExist(ctx context.Context, fqfield string, value interface{}) error {
	if err := dp.Get(ctx, fqfield, value); err != nil {
		var errDoesNotExist doesNotExistError
		if !errors.As(err, &errDoesNotExist) {
			return err
		}
	}
	return nil
}

// Exists tells, if a fqfield exist.
//
// If an error happens, it returns false.
func (dp DataProvider) Exists(ctx context.Context, fqfield string) (bool, error) {
	fields, err := dp.externalGet(ctx, fqfield)
	if err != nil {
		return false, fmt.Errorf("getting fqfield: %w", err)
	}

	return fields[0] != nil, nil
}
