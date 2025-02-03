package handle

import (
	"net/http"
)

type ContextFunc func(c *Context)

func (cf ContextFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	cf(c)
}
