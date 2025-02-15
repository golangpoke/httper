package httper_test

import (
	"github.com/golangpoke/httper"
	"github.com/golangpoke/httper/handle"
	"github.com/golangpoke/httper/middle"
	"github.com/golangpoke/httper/result"
	"github.com/golangpoke/nlog"
	"testing"
)

type Data struct {
	Hello string `json:"hello" validate:"required"`
}

func TestServeMux_Run(t *testing.T) {
	defer nlog.Recovery()
	nlog.SetDefault(nlog.LvlDebug)
	mux := httper.NewServeMux()
	mux.Use(middle.CORS(), middle.Logger(), middle.Recovery())
	mux.POST("/", handle.ResultHandle(func(c *handle.Context) result.Result {
		var data Data
		err := c.BindValidJson(&data, "Hello")
		if err != nil {
			return result.ErrBadRequest.Wrap(err)
		}
		// err := c.BindJSON(&data)
		// if err != nil {
		// 	return result.ErrBadRequest.Wrap(err)
		// }
		// err = validator.New().StructPartial(data, "Hello")
		// if err != nil {
		// 	return result.ErrBadRequest.Wrap(err)
		// }
		return result.Map{
			"data": data,
		}
	}))
	for r, f := range mux.Routes() {
		t.Log("register router:", r, f)
	}
	t.Fatal(mux.Start(":8000"))
}
