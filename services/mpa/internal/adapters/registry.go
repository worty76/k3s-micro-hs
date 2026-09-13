package adapters

type Registry struct {
	adapters map[Protocol]Runnable
}

func NewRegistry(adapters map[Protocol]Runnable) *Registry {
	return &Registry{
		adapters: adapters,
	}
}

func (r *Registry) Get(protocol Protocol) (Runnable, bool) {
	adapter, ok := r.adapters[protocol]
	return adapter, ok
}

func (r *Registry) All() []Runnable {
	adapters := make([]Runnable, 0, len(r.adapters))

	for _, adapter := range r.adapters {
		adapters = append(adapters, adapter)
	}

	return adapters
}
