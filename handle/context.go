package handle

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/golangpoke/nlog"
	"io"
	"net/http"
)

type Context struct {
	writer     http.ResponseWriter
	request    *http.Request
	statusCode int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{writer: w, request: r}
}

func (c *Context) JSON(status int, data any) error {
	c.WriteHeader(status)
	c.writer.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c.writer).Encode(data)
}

func (c *Context) Error(code int, err error) error {
	c.statusCode = code
	return err
}

func (c *Context) PathValue(name string) string {
	return c.request.PathValue(name)
}

func (c *Context) URLPath() string {
	return c.request.URL.Path
}

func (c *Context) Method() string {
	return c.request.Method
}

func (c *Context) BindJSON(data any) error {
	body, err := io.ReadAll(c.request.Body)
	if err != nil {
		return nlog.Catch(err)
	}
	err = json.Unmarshal(body, data)
	if err != nil {
		return nlog.Catch(err)
	}
	return nil
}

func (c *Context) BindValidJson(data any, fields ...string) error {
	err := c.BindJSON(data)
	if err != nil {
		return nlog.Catch(err)
	}
	err = validator.New().StructPartial(data, fields...)
	if err != nil {
		return nlog.Catch(err)
	}
	return nil
}

func (c *Context) StatusCode() int {
	return c.statusCode
}

func (c *Context) Header() http.Header {
	return c.writer.Header()
}

func (c *Context) Write(bytes []byte) (int, error) {
	return c.writer.Write(bytes)
}

func (c *Context) WriteHeader(statusCode int) {
	c.statusCode = statusCode
	c.writer.WriteHeader(statusCode)
}

func (c *Context) SetHeader(name, value string) {
	c.writer.Header().Set(name, value)
}

// ProxyHandle handle next http.handler,return statusCode
func (c *Context) ProxyHandle(next http.Handler) {
	w := &Context{
		writer: c.writer,
	}
	c.writer = w
	next.ServeHTTP(w, c.request)
	c.statusCode = w.statusCode
}
