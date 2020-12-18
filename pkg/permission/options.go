package permission

import "github.com/OpenSlides/openslides-permission-service/internal/collection"

// Option is an optional argument for permission.New()
type Option func(*Permission)

// WithConnecters initializes a Permission Service with specific connecters. Per
// default, the OpenSlides collections are used.
func WithConnecters(cons []collection.Connecter) Option {
	return func(p *Permission) {
		p.connecters = cons
	}
}
