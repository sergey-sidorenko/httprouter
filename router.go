package netrouter

import (
	"net/http"
)

type Router struct {
	rules []*Rule
}

func (r *Router) AddRule(rule *Rule) {
	r.rules = append(r.rules, rule)
}
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var rule *Rule
	for _, curRule := range r.rules {
		if rule.Match(req) {
			rule = curRule
			break
		}
	}
	rule.Handle(w, req)
}
