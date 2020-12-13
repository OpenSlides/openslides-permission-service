package permission

// Option is an optional argument for permission.New()
type Option func(*Permission)

// WithCollections initializes a Permission Service with specific collections.
// Per default, the OpenSlides collections are used.
func WithCollections(coll map[string]Collection) Option {
	return func(p *Permission) {
		p.coll = coll
	}
}
