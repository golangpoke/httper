package handle_test

import (
	"github.com/golangpoke/httper"
	"github.com/golangpoke/httper/handle"
	"github.com/golangpoke/httper/middle"
	"github.com/golangpoke/httper/result"
	"github.com/golangpoke/nlog"
	"net/http"
	"os"
	"testing"
)

func TestUnitHandle(t *testing.T) {
	nlog.Use()
	mux := httper.NewServeMux()
	mux.Use(middle.CORS(), middle.Recovery(), middle.Logger())
	mux.POST("/none", handle.ContextFunc(func(c *handle.Context) {
		_ = c.JSON(http.StatusInternalServerError, "ok")
	}))
	mux.POST("/http", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	mux.POST("/error", handle.ErrorFunc(func(c *handle.Context) error {
		err := os.ErrNotExist
		return c.Error(http.StatusInternalServerError, err)
	}))
	mux.POST("/result", handle.ResultHandle(func(c *handle.Context) result.Result {
		data := struct {
			Hello string `json:"hello"`
		}{}
		err := c.BindJSON(&data)
		if err != nil {
			return result.ErrBadRequest.Wrap(err)
		}
		return result.Map{
			"data": data,
		}
	}))
	for r, f := range mux.Routes() {
		t.Log("register router:", r, f)
	}
	t.Fatal(mux.Start(":8000"))
}
