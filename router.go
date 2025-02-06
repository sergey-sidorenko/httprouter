package httprouter

import (
	"net/http"
)

type Router struct {
	rules map[string]*Rule
}

func (r *Router) AddRule(rule *Rule) {
	key := rule.pattern + string(rule.method)
	r.rules[key] = rule
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
	return &Router{make(map[string]*Rule, 0)}
}
