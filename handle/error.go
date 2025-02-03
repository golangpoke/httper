package handle

import (
	"net/http"
)

type ErrorFunc func(c *Context) error

func (ef ErrorFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	err := ef(c)
	UnitErrorHandle(c, err)
}

var UnitErrorHandle = func(c *Context, err error) {
	if err == nil {
		return
	}
	http.Error(c, err.Error(), c.StatusCode())
}
