package conch

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Response 	http.ResponseWriter
	Request 	*http.Request
	App			App
}

type Controller map[string]func(ctx Context)
type Param map[string]interface{}

func (ctx Context) Get(name string) string {
	if ctx.Request.Method == "POST" {
		return ctx.Request.PostFormValue(name)
	}
	return ctx.Request.FormValue(name)
}

func (ctx Context) Json(data ...interface{}) {
	str, err := json.Marshal(data)
	if err != nil {
		_, _ = ctx.Response.Write([]byte(err.Error()))
		return
	}
	ctx.Response.Header().Set("Content-Type", "application/json")
	_, _ = ctx.Response.Write(str)
}

func (ctx Context) Log(format string, arg ...interface{}) {
	Log{Path: ctx.App.LogPath}.Out(format, arg...)
}