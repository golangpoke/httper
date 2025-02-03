package handle

import (
	"github.com/golangpoke/httper/result"
	"net/http"
)

type ResultHandle func(c *Context) result.Result

func (rh ResultHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	rs := rh(c)
	UnitResultHandle(c, rs)
}

var UnitResultHandle = func(c *Context, rs result.Result) {
	if rs == nil {
		return
	}
	status := http.StatusOK
	m := make(map[string]any)
	if rs.Error() != nil || rs.Code() != result.Success {
		status = http.StatusInternalServerError
		m["error"] = rs.Error().Error()
	}
	if msg := result.MapCodeMessage[rs.Code()]; msg != "" {
		m["message"] = msg
	}
	if rs.Data() != nil {
		m["data"] = rs.Data()
	}
	m["code"] = rs.Code()
	err := c.JSON(status, m)
	if err != nil {
		http.Error(c, err.Error(), http.StatusInternalServerError)
	}
}
