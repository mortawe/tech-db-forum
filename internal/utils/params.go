package utils

import (
	"github.com/mortawe/tech-db-forum/internal/models"
	"github.com/valyala/fasthttp"
	"strconv"
)

func MustGetParams(ctx *fasthttp.RequestCtx) *models.Params {
	params := &models.Params{}
	params.Since = string(ctx.FormValue("since"))
	params.Sort = string(ctx.FormValue("sort"))
	if string(ctx.FormValue("desc")) == "true" {
		params.Desc = true
	} else {
		params.Desc = false
	}
	params.Limit, _ = strconv.Atoi(string(ctx.FormValue("limit")))
	return params
}
