package routing

import (
	"container/list"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

const maxAmountOfRules int = 1024

type MiddlewareFunc func(req *Request, resp *Response)

// Router - http network router
type Router struct {
	handlers                   *list.List
	handlersMux                sync.Mutex
	handlerCache               map[string]*list.Element
	beforeHandleMiddlewares    []MiddlewareFunc
	beforeHandleMiddlewaresMux sync.Mutex
	inCT                       ContentType
	outCT                      ContentType
}

// NewRouter - created new Router
func NewRouter() *Router {
	return &Router{list.New(), sync.Mutex{}, make(map[string]*list.Element), make([]MiddlewareFunc, 0), sync.Mutex{}, ContentTypeJson, ContentTypeJson}
}

// SetInContentType - sets incoming content type
func (r *Router) SetInContentType(ct ContentType) {
	r.inCT = ct
}

// SetOutContentType - sets outcoming content type
func (r *Router) SetOutContentType(ct ContentType) {
	r.outCT = ct
}

// findHandlerElement - finds element in rules list
func (r *Router) findHandlerElement(pattern string, method string) *list.Element {
	var el *list.Element
	var route RouteHandler
	var found bool
	method = strings.ToLower(method)
	for curEl := r.handlers.Front(); !found && curEl != nil; curEl = curEl.Next() {
		route = curEl.Value.(RouteHandler)
		if route.Pattern() == pattern && route.Method() == method {
			el = curEl
			found = true
		}
	}
	return el
}

// AddRoute - adds a route
func (r *Router) AddRoute(route RouteHandler) {
	var existEl *list.Element = r.findHandlerElement(route.Pattern(), route.Method())
	r.handlersMux.Lock()
	defer r.handlersMux.Unlock()
	if existEl != nil {
		existEl.Value = route
		r.handlers.MoveToFront(existEl)
	} else {
		r.handlers.PushFront(route)
	}
	if r.handlers.Len() > maxAmountOfRules {
		r.handlers.Remove(r.handlers.Back())
	}
}

// RemoveRoute - removes a route
func (r *Router) RemoveRoute(pattern string, method string) {
	var el *list.Element = r.findHandlerElement(pattern, method)
	if el != nil {
		r.handlers.Remove(el)
	}
}

// Handle - find first route whose rule matches the Request
func (r *Router) Route(req *Request, resp *Response) (RouteHandler, error) {
	var routeHandler RouteHandler
	var params PathParams
	err := fmt.Errorf("%w", MatchErrorNoMatch)
	for el := r.handlers.Front(); el != nil; el = el.Next() {
		routeHandler = el.Value.(RouteHandler)
		params, err = routeHandler.Match(req)
		if err == nil {
			req.params = params
			return routeHandler, nil
		}
	}
	return nil, err
}

// ServeHTTP - makes Router acts as multiplexer
func (r *Router) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	req := &Request{request, nil}
	resp := &Response{w}
	route, err := r.Route(req, resp)
	if route != nil {
		for _, middleware := range r.beforeHandleMiddlewares {
			middleware(req, resp)
		}
		route.Run(req, resp)
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

func (r *Router) BeforeHandle(f MiddlewareFunc) {
	r.beforeHandleMiddlewaresMux.Lock()
	defer r.beforeHandleMiddlewaresMux.Unlock()
	r.beforeHandleMiddlewares = append(r.beforeHandleMiddlewares, f)
}
