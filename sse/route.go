package sse

import "github.com/sergey-sidorenko/httprouter/routing"

// HandlerFunc - function-handler to handle http-request
type HandlerFunc func(req *Request, resp MessageWriter)

// Route - rule for url-to-action routing
type Route struct {
	Rule    *routing.Rule
	Handler HandlerFunc
}

// Run - runs the handler to handle http request
func (r Route) Run(req *routing.Request, resp *routing.Response) {
	r.Handler(&Request{req}, &Response{resp})
}

// Match - matches the handler against the http request
func (r Route) Match(req *routing.Request) (routing.PathParams, error) {
	return r.Rule.Match(req)
}

// Pattern - returns pattern
func (r Route) Pattern() string {
	return r.Rule.Pattern()
}

// Method - returns http request method that corresponds to the route
func (r Route) Method() string {
	return r.Rule.Method()
}
