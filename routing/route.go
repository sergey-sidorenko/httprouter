package routing

type RouteHandler interface {
	Match(req *Request) (PathParams, error)
	Run(req *Request, resp *Response)
	Pattern() string
	Method() string
}

// HandlerFunc - function-handler to handle http-request
type HandlerFunc func(req *Request, resp *Response)

// Route - rule for url-to-action routing
type Route struct {
	Rule    *Rule
	Handler HandlerFunc
}

// Run - runs the handler to handle http request
func (r Route) Run(req *Request, resp *Response) {
	r.Handler(req, resp)
}

// Match - matches the handler against the http request
func (r Route) Match(req *Request) (PathParams, error) {
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
