package sse

import "github.com/sergey-sidorenko/httprouter/routing"

// Request - wrapper for SSE Request
type Request struct {
	req *routing.Request
}

func (r *Request) PathParam(key string) (string, error) {
	return r.req.PathParam(key)
}
