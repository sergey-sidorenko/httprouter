package routing

import "net/http"

// Response - wrapper for http Response
type Response struct {
	http.ResponseWriter
}
