package middle

import (
	"fmt"
	"github.com/golangpoke/httper"
	"github.com/golangpoke/httper/handle"
	"github.com/golangpoke/nlog"
	"net/http"
)

// CORS allow all origin/methods/headers, may not safe
func CORS() httper.Middleware {
	return func(next http.Handler) http.Handler {
		return handle.ContextFunc(func(c *handle.Context) {
			c.SetHeader("Access-Control-Allow-Origin", "*")
			c.SetHeader("Access-Control-Allow-Methods", "*")
			c.SetHeader("Access-Control-Allow-Headers", "*")
			if c.Method() == http.MethodOptions {
				c.WriteHeader(http.StatusNoContent)
				return
			}
			c.ProxyHandle(next)
		})
	}
}

// Recovery use nlog recovery
func Recovery() httper.Middleware {
	return func(next http.Handler) http.Handler {
		return handle.ContextFunc(func(c *handle.Context) {
			defer nlog.Recovery()
			c.ProxyHandle(next)
		})
	}
}

func Logger() httper.Middleware {
	return func(next http.Handler) http.Handler {
		return handle.ContextFunc(func(c *handle.Context) {
			c.ProxyHandle(next)
			s := fmt.Sprintf("method:%s status:%d url:%s", c.Method(), c.StatusCode(), c.URLPath())
			switch c.StatusCode() {
			case http.StatusOK:
				nlog.Option(nlog.NoSource()).INFOf(s)
			default:
				nlog.Option(nlog.NoSource()).WARNf(s)
			}
		})
	}
}
