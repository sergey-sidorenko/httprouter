package sse

import (
	"net/http"

	"github.com/sergey-sidorenko/httprouter/routing"
)

// Server - SSE server
type Server struct {
	r   *routing.Router
	srv *http.Server
}

// NewServer - created new SSE server
func NewServer(address string, router *routing.Router) *Server {
	router.BeforeHandle(func(req *routing.Request, resp *routing.Response) {
		resp.Header().Set("X-Accel-Buffering", "no")
		resp.Header().Set("Content-Type", "text/event-stream")
		resp.Header().Set("Cache-Control", "no-cache")
		resp.Header().Set("Access-Control-Allow-Origin", "http://web.local")
	})
	httpServer := &http.Server{Addr: address, Handler: router}
	sseServer := &Server{r: router, srv: httpServer}
	return sseServer
}

// Run - starts the SSE server
func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}

// SetRouter - sets Router
func (s *Server) SetRouter(router *routing.Router) {
	s.srv.Handler = router
}
