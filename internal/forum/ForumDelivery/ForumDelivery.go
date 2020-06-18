package ForumDelivery

import (
	"encoding/json"
	"github.com/fasthttp/router"
	"github.com/mortawe/tech-db-forum/internal/forum"
	"github.com/mortawe/tech-db-forum/internal/models"
	"github.com/mortawe/tech-db-forum/internal/user"
	"github.com/mortawe/tech-db-forum/internal/utils"
	"github.com/valyala/fasthttp"
)

type ForumManager struct {
	fUC forum.IForumUC
	uUC user.IUserUC
}

func NewForumManager(f forum.IForumUC, u user.IUserUC) *ForumManager {
	return &ForumManager{fUC: f, uUC: u}
}

func (m *ForumManager) InitRoutes(r *router.Router) {
	r.POST("/api/forum/create", m.CreateForum)
	r.GET("/api/forum/{slug}/details", m.Details)
	r.GET("/api/forum/{slug}/users", m.GetUsersByForum)
}

func (m *ForumManager) CreateForum(ctx *fasthttp.RequestCtx) {
	forum := &models.Forum{}
	if err  := forum.UnmarshalJSON(ctx.PostBody()); err != nil {
		utils.Send(400, ctx, utils.MustMarshalError(err))
		return
	}
	var err error
	forum.User, err = m.uUC.SelectNicknameWithCase(forum.User)
	if err != nil {
		utils.Send(404, ctx, utils.MustMarshalError(err))
		return
	}
	err = m.fUC.Create(forum)
	
	switch err {
	case models.ErrNotExists:
		utils.Send(404, ctx, utils.MustMarshalError(err))
	case models.ErrConflict:
		forumInBase, _ := m.fUC.SelectBySlug(forum.Slug)
		resp, _ := forumInBase.MarshalJSON()
		utils.Send(409, ctx, resp)
	case nil:
		resp, _ := json.Marshal(forum)
		utils.Send(201, ctx, resp)
	default:
		utils.Send(500, ctx, utils.MustMarshalError(err))
	}
}

func (m *ForumManager) Details(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	forum, err := m.fUC.SelectBySlug(slug)
	switch err {
	case nil:
		resp, _ := forum.MarshalJSON()
		utils.Send(200, ctx, resp)
	default:
		utils.Send(404, ctx, utils.MustMarshalError(err))
	}
}

func (m *ForumManager) GetUsersByForum(ctx *fasthttp.RequestCtx) {
	params := utils.MustGetParams(ctx)
	slug := ctx.UserValue("slug").(string)
	var err error
	if slug, err = m.fUC.SelectForumWithCase(slug); err != nil {
		utils.Send(404, ctx, utils.MustMarshalError(err))
		return
	}
	users, err := m.fUC.GetUsersByForum(slug, params.Desc, params.Since, params.Limit)

	switch err {
	case nil:
		resp, _ := json.Marshal(users)
		utils.Send(200, ctx, resp)
	default:
		utils.Send(404, ctx, utils.MustMarshalError(err))
	}

}