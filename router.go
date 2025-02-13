package httprouter

import (
	"container/list"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

const maxAmountOfRules int = 1024

// Router - http network router
type Router struct {
	rules      *list.List
	rulesCache map[string]*Rule
	rlmutex    sync.Mutex
}

// findRuleElement - finds element in rules list
func (r *Router) findRuleElement(pattern string, method string) *list.Element {
	var el *list.Element
	var rule *Rule
	var found bool
	method = strings.ToLower(method)
	for curEl := r.rules.Front(); !found && curEl != nil; curEl = curEl.Next() {
		rule = curEl.Value.(*Rule)
		if rule.pattern == pattern && rule.method == method {
			el = curEl
			found = true
		}
	}
	return el
}

// AddRule - adds a rule
func (r *Router) AddRule(rule *Rule) {
	var existEl *list.Element = r.findRuleElement(rule.pattern, rule.method)
	r.rlmutex.Lock()
	defer r.rlmutex.Unlock()
	if existEl != nil {
		existEl.Value = rule
		r.rules.MoveToFront(existEl)
	} else {
		r.rules.PushFront(rule)
	}
	if r.rules.Len() > maxAmountOfRules {
		r.rules.Remove(r.rules.Back())
	}
}

// RemoveRule - removes a rule
func (r *Router) RemoveRule(pattern string, method string) {
	var el *list.Element = r.findRuleElement(pattern, method)
	if el != nil {
		r.rules.Remove(el)
	}
}

// MatchedRule - find first rule that match the Request
func (r *Router) MatchedRule(w http.ResponseWriter, req *http.Request) (*Rule, RuleContext, error) {
	var ctx RuleContext
	var rule *Rule
	err := fmt.Errorf("%w", MatchErrorNoMatch)
	for el := r.rules.Front(); el != nil; el = el.Next() {
		rule = el.Value.(*Rule)
		ctx, err = rule.Match(req)
		if err == nil {
			return rule, ctx, nil
		}
	}
	return nil, ctx, err
}

// ServeHTTP - makes Router acts as multiplexer
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rule, ctx, err := r.MatchedRule(w, req)
	if rule != nil {
		rule.Handle(w, req, ctx)
	} else {
		r.HandleError(w, err)
	}
}

// HandleError - handles error fired up during finding a matched rule
func (r *Router) HandleError(w http.ResponseWriter, err error) {
	var status int
	var rme *RuleMatchError
	if errors.As(err, &rme) {
		switch rme.Code {
		case MatchErrorCodeWrongMethod:
			status = http.StatusMethodNotAllowed
		case MatchErrorCodeNoMatch:
			status = http.StatusNotFound
		case MatchErrorCodeWrongQueryParam:
			status = http.StatusBadRequest
		}
	} else {
		status = http.StatusInternalServerError
	}
	w.WriteHeader(status)
}

// NewRouter - created new Router
func NewRouter() *Router {
	return &Router{list.New(), make(map[string]*Rule), sync.Mutex{}}
}
