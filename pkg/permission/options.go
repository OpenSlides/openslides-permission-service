package permission

import "github.com/OpenSlides/openslides-permission-service/internal/types"

// Option is an optional argument for permission.New()
type Option func(*Permission)

// WithConnecters initializes a Permission Service with specific connecters. Per
// default, the OpenSlides collections are used.
func WithConnecters(cons []types.Connecter) Option {
	return func(p *Permission) {
		p.connecters = cons
	}
}
