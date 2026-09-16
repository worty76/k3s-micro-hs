package integrations

type Registry struct {
	integrations map[string]Integration
}

func NewRegistry() *Registry {
	return &Registry{
		integrations: make(map[string]Integration),
	}
}

func (r *Registry) Register(integration Integration) {
	r.integrations[integration.Name()] = integration
}

func (r *Registry) Get(name string) (Integration, bool) {
	integration, exists := r.integrations[name]

	if !exists {
		return nil, false
	}

	return integration, exists
}
