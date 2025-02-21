package sse

import (
	"net/http"

	"github.com/sergey-sidorenko/httprouter/routing"
)

type MessageWriter interface {
	WriteMessage(m *Message) error
}

// Response - wrapper for SSE Response
type Response struct {
	*routing.Response
}

func (resp *Response) WriteMessage(m *Message) error {
	_, err := resp.Write(m.Bytes())
	if err != nil {
		return err
	}
	if f, ok := resp.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}
