package mqtt

import (
	"strings"

	"github.com/worty76/k3s-micro-hs/services/mpa/internal/adapters/integrations"
)

type Router struct {
	routes []route
}

type route struct {
	pattern []string
	target  integrations.Integration
}

func (r *Router) Register(pattern string, target integrations.Integration) {
	r.routes = append(r.routes, route{
		pattern: strings.Split(pattern, "/"),
		target:  target,
	})
}

func (r *Router) Match(topic string) (integrations.Integration, bool) {
	levels := strings.Split(topic, "/")
	for _, rt := range r.routes {
		if match(rt.pattern, levels) {
			return rt.target, true
		}
	}
	return nil, false
}

func match(pattern, topic []string) bool {
	for i, p := range pattern {
		if p == "#" {
			return true
		}
		if i >= len(topic) {
			return false
		}
		if p != "+" && p != topic[i] {
			return false
		}
	}
	return len(pattern) == len(topic)
}
