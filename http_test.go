package httper_test

import (
	"github.com/golangpoke/httper"
	"github.com/golangpoke/httper/handle"
	"net/http"
	"testing"
)

func TestServeMux_Run(t *testing.T) {
	mux := httper.NewServeMux()
	mux.POST("/", handle.ContextFunc(func(c *handle.Context) {
		c.JSON(http.StatusInternalServerError, "ok")
	}))
	for r, f := range mux.Routes() {
		t.Log("register router:", r, f)
	}
	t.Fatal(mux.Start(":8000"))
}
