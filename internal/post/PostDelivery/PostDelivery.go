package PostDelivery

import (
	"bytes"
	"encoding/json"
	"github.com/fasthttp/router"
	"github.com/mortawe/tech-db-forum/internal/forum"
	"github.com/mortawe/tech-db-forum/internal/models"
	"github.com/mortawe/tech-db-forum/internal/post"
	"github.com/mortawe/tech-db-forum/internal/thread"
	"github.com/mortawe/tech-db-forum/internal/user"
	"github.com/mortawe/tech-db-forum/internal/utils"
	"github.com/valyala/fasthttp"
	"strconv"
)

type PostManager struct {
	fUC forum.IForumUC
	pUC post.IPostUC
	tUC thread.IThreadUC
	uUC user.IUserUC
}

func NewForumManager(f forum.IForumUC, p post.IPostUC, t thread.IThreadUC, u user.IUserUC) *PostManager {
	return &PostManager{fUC: f, pUC: p, tUC: t, uUC: u}
}

func (m *PostManager) InitRoutes(r *router.Router) {
	r.POST("/api/thread/{slugOrID}/create", m.Create)
	r.POST("/api/post/{id}/details", m.Update)
	r.GET("/api/post/{id}/details", m.GetByID)
	r.GET("/api/thread/{slugOrID}/posts", m.GetPosts)
}

func (m *PostManager) Create(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slugOrID").(string)
	id, forum, err := m.tUC.GetIDForumBySlugOrID(slugOrID)
	if err != nil {
		utils.Send(404, ctx, utils.MustMarshalError(err))
		return
	}
	posts := &[]*models.Post{}
	if err := json.Unmarshal(ctx.PostBody(), posts); err != nil {
		utils.Send(400, ctx, utils.MustMarshalError(err))
		return
	}
	err = m.pUC.InsertPost(*posts, forum, id)
	switch err {
	case nil:
		resp, _ := json.Marshal(posts)
		utils.Send(201,ctx, resp)
	case models.ErrConflict:
		utils.Send(409, ctx, utils.MustMarshalError(err))
	default:
		utils.Send(404, ctx, utils.MustMarshalError(err))
	}
}

func (m *PostManager) Update(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, _ := strconv.Atoi(idStr)
	post := &models.Post{}
	if err := json.Unmarshal(ctx.PostBody(), post); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(`{"message": "` + "unmarshal not ok : " + err.Error() + `"}`))
		return
	}
	post.ID = id
	err := m.pUC.Update(post)
	switch err {
	case nil:
		resp, _ := post.MarshalJSON()
		utils.Send(200, ctx, resp)
	default:
		utils.Send(404, ctx, utils.MustMarshalError(err))
	}
}

func (m *PostManager) GetByID(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	related := ctx.QueryArgs().Peek("related")
	idStr := ctx.UserValue("id").(string)
	id, _ := strconv.Atoi(idStr)
	details := &models.PostDetails{}
	var err error
	details.Post, err = m.pUC.SelectPostByID(id)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "`  + err.Error() + `"}`))
		return
	}
	user, err  := m.uUC.SelectByNickname(details.Post.Author)
	details.User = &user
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "`  + err.Error() + `"}`))
		return
	}
	details.Thread, err = m.tUC.SelectByID(details.Post.Thread)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "`  + err.Error() + `"}`))
		return
	}
	details.Forum, err  = m.fUC.SelectBySlug(details.Thread.Forum)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "`  + err.Error() + `"}`))
		return
	}

	if !bytes.Contains(related, []byte("user")) {
		details.User = nil
	}
	if !bytes.Contains(related, []byte("forum")) {
		details.Forum = nil
	}
	if !bytes.Contains(related, []byte("thread")) {
		details.Thread = nil
	}

	resp, _ := details.MarshalJSON()
	ctx.Write(resp)
	ctx.SetStatusCode(200)
}

func (m *PostManager) GetPosts(ctx *fasthttp.RequestCtx) {
	params := utils.MustGetParams(ctx)
	slug := ctx.UserValue("slugOrID").(string)
	thread, err := m.tUC.SelectBySlugOrID(slug)
	if err != nil {
		utils.Send(404, ctx, utils.MustMarshalError(err))
		return
	}
	posts, err := m.pUC.GetPosts(thread.ID, params.Desc, params.Since, params.Limit, params.Sort)
	switch err {
	case nil:
		resp, _ := json.Marshal(posts)
		utils.Send(200, ctx, resp)
	default:
		utils.Send(404, ctx, utils.MustMarshalError(err))
	}
}