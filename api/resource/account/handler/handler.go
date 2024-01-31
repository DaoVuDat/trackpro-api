package accounthandler

import "trackpro/util/ctx"

type Handler struct {
}

func New(app *ctx.Application) *Handler {
	return &Handler{}
}
