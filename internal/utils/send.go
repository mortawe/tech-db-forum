package utils

import (
	"github.com/mortawe/tech-db-forum/internal/models"
	"github.com/valyala/fasthttp"
)

func Send(status int, ctx *fasthttp.RequestCtx, resp []byte) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(status)
	ctx.Write(resp)
}

func MustMarshalError(err error) []byte{
	m := models.Msg{Message: err.Error()}
	resp, _ := m.MarshalJSON()
	return resp
}
