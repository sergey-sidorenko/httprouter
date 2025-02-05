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
		if curRule.Match(req) {
			rule = curRule
			break
		}
	}
	if rule != nil {
		rule.Handle(w, req)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Requested page not found!"))
	}
}
func NewRouter() *Router {
	return &Router{make([]*Rule, 0)}
}
