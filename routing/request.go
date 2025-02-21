package routing

import (
	"errors"
	"net/http"
)

// Request - wrapper for http Request
type Request struct {
	*http.Request
	params PathParams
}

func (r *Request) PathParam(key string) (string, error) {
	var err error = errors.New("no param with such key found")
	if r.params != nil {
		return "", err
	}
	value, ok := r.params[key]
	if ok {
		return value, nil
	}
	return "", err
}
