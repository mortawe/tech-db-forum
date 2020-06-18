package ThreadDelivery

import (
	"encoding/json"
	"errors"
	"github.com/fasthttp/router"
	"github.com/mortawe/tech-db-forum/internal/forum"
	"github.com/mortawe/tech-db-forum/internal/models"
	"github.com/mortawe/tech-db-forum/internal/thread"
	"github.com/mortawe/tech-db-forum/internal/user"
	"github.com/mortawe/tech-db-forum/internal/utils"
	"github.com/valyala/fasthttp"
)

type ThreadManager struct {
	fUC forum.IForumUC
	uUC user.IUserUC
	tUC thread.IThreadUC
}

func NewThreadManager(fUC forum.IForumUC, uUC user.IUserUC, tUC thread.IThreadUC) *ThreadManager {
	return &ThreadManager{
		fUC: fUC,
		uUC: uUC,
		tUC: tUC,
	}
}

func (m *ThreadManager) InitRoutes(r *router.Router) {
	r.GET("/api/forum/{slug}/threads", m.GetThreadsByForum)
	r.GET("/api/thread/{slugOrID}/details", m.Details)
	r.POST("/api/forum/{slug}/create", m.CreateThread)
	r.POST("/api/thread/{slugOrID}/details", m.Update)
	r.POST("/api/thread/{slugOrID}/vote", m.Vote)
}

func (m *ThreadManager) CreateThread(ctx *fasthttp.RequestCtx) {
	thread := &models.Thread{}
	if err := thread.UnmarshalJSON(ctx.PostBody()); err != nil {
		utils.Send(400, ctx, utils.MustMarshalError(err))
		return
	}
	thread.Forum = ctx.UserValue("slug").(string)
	var err error
	thread.Forum, err = m.fUC.SelectForumWithCase(thread.Forum)
	if err != nil {
		utils.Send(404, ctx, utils.MustMarshalError(err))
		return
	}
	nickname, err := m.uUC.SelectNicknameWithCase(thread.Author)
	if err != nil {
		utils.Send(404, ctx, utils.MustMarshalError(err))
		return
	}
	thread.Author = nickname
	threadInBase := models.Thread{}
	if thread.Slug != "" {
		threadInBase, err := m.tUC.SelectThreadBySlug(thread.Slug)
		if err == nil {
			resp, _ := threadInBase.MarshalJSON()
			utils.Send(409, ctx, resp)
			return
		}
	}

	err = m.tUC.InsertThread(thread)
	switch err {
	case nil:
		resp, _ := thread.MarshalJSON()
		utils.Send(201, ctx, resp)
	case models.ErrConflict:
		resp, _ := threadInBase.MarshalJSON()
		utils.Send(409, ctx, resp)
	default:
		utils.Send(500, ctx, utils.MustMarshalError(err))
	}
}

func (m *ThreadManager) GetThreadsByForum(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	params := utils.MustGetParams(ctx)

	if _, err := m.fUC.SelectForumWithCase(slug); err != nil {
		utils.Send(404, ctx,
			utils.MustMarshalError(errors.New("forum doesn't exists")))
		return
	}

	threads, err := m.tUC.SelectThreadsByForum(slug, params.Limit, params.Since, params.Desc)
	switch err {
	case nil:
		resp, _ := json.Marshal(threads)
		utils.Send(200, ctx, resp)
	default:
		utils.Send(404, ctx, utils.MustMarshalError(err))
	}
}

func (m *ThreadManager) Details(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slugOrID").(string)
	thread, err := m.tUC.SelectBySlugOrID(slug)
	switch err {
	case nil:
		resp, _ := thread.MarshalJSON()
		utils.Send(200, ctx, resp)
	default:
		utils.Send(404, ctx, utils.MustMarshalError(err))
	}
}

func (m *ThreadManager) Update(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slugOrID").(string)
	thread := &models.Thread{}
	if err := thread.UnmarshalJSON(ctx.PostBody()); err != nil {
		utils.Send(400, ctx, utils.MustMarshalError(err))
	}
	err := m.tUC.UpdateBySlugOrID(slug, thread)
	switch err {
	case nil:
		resp, _ := thread.MarshalJSON()
		utils.Send(200, ctx, resp)
	default:
		utils.Send(404, ctx, utils.MustMarshalError(err))
	}
}

func (m *ThreadManager) Vote(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slugOrID").(string)

	voice := &models.Vote{}
	if err := voice.UnmarshalJSON(ctx.PostBody()); err != nil {
		utils.Send(400, ctx, utils.MustMarshalError(err))
	}
	thread, err := m.tUC.Vote(*voice, slug)
	switch err {
	case nil:
		resp, _ := thread.MarshalJSON()
		utils.Send(200, ctx, resp)
	default:
		utils.Send(404, ctx, utils.MustMarshalError(err))
	}
}
