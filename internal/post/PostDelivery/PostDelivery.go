package PostDelivery

import (
	"bytes"
	"encoding/json"
	"github.com/fasthttp/router"
	"github.com/mortawe/tech-db-forum/internal/forum"
	"github.com/mortawe/tech-db-forum/internal/models"
	"github.com/mortawe/tech-db-forum/internal/post"
	"github.com/mortawe/tech-db-forum/internal/post/PostUC"
	"github.com/mortawe/tech-db-forum/internal/thread"
	"github.com/mortawe/tech-db-forum/internal/user"
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
	ctx.SetContentType("application/json")

	slugOrID := ctx.UserValue("slugOrID").(string)
	thread, err := m.tUC.SelectBySlugOrID(slugOrID)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "` + err.Error() + " slug or id : " + slugOrID +"\"}"))
		return
	}
	posts := &[]*models.Post{}
	if err := json.Unmarshal(ctx.PostBody(), posts); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(`{"message": "` + "unmarshal not ok : " + err.Error() + `"}`))
		return
	}

	for _, p := range *posts {
		if _, err := m.uUC.SelectByNickname(p.Author); err != nil {
			ctx.SetStatusCode(404)
			ctx.Write([]byte(`{"message": "` + "user error : " + err.Error() + `"}`))
			return
		}
	}
	err = m.pUC.InsertPost(*posts, thread.Forum, thread.ID)
	if err != nil {
		if err == PostUC.ParentErr {
			ctx.SetStatusCode(409)
			ctx.Write([]byte(`{"message": "` + "parent error : "+ `"}`))
			return
		}
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(`{"message": "` + "insert not ok : " + `"}`))
		return
	}
	resp, _ := json.Marshal(posts)
	ctx.Write(resp)
	ctx.SetStatusCode(201)
}

func (m *PostManager) Update(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	idStr := ctx.UserValue("id").(string)
	id, _ := strconv.Atoi(idStr)
	postInDB, err := m.pUC.SelectPostByID(id)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "` + "can't find this post : " + err.Error() + `"}`))
		return
	}
	post := &models.Post{}
	if err := json.Unmarshal(ctx.PostBody(), post); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(`{"message": "` + "unmarshal not ok : " + err.Error() + `"}`))
		return
	}
	post.ID = id
	if postInDB.Message == post.Message || post.Message == "" {
		resp, _ := json.Marshal(postInDB)
		ctx.Write(resp)
		ctx.SetStatusCode(200)
		return
	}
	if err := m.pUC.Update(post); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(`{"message": "` + "update not ok : " + err.Error() + `"}`))
		return
	}
	resp, _ := json.Marshal(post)
	ctx.Write(resp)
	ctx.SetStatusCode(200)
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

	resp, _ := json.Marshal(details)
	ctx.Write(resp)
	ctx.SetStatusCode(200)
}

func (m *PostManager) GetPosts(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	params := &models.GetThreadsParams{}
	params.Since = string(ctx.FormValue("since"))
	if string(ctx.FormValue("desc")) == "true" {
		params.Desc = true
	} else {
		params.Desc = false
	}
	params.Limit, _ = strconv.Atoi(string(ctx.FormValue("limit")))
	params.Sort = string(ctx.FormValue("sort"))
	slug := ctx.UserValue("slugOrID").(string)
	thread, err := m.tUC.SelectBySlugOrID(slug)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "`  + err.Error() + `"}`))
		return
	}

	posts, err := m.pUC.GetPosts(thread.ID, params.Desc, params.Since, params.Limit, params.Sort)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "`  + err.Error() + `"}`))
		return
	}

	resp, _ := json.Marshal(posts)
	ctx.Write(resp)
	ctx.SetStatusCode(200)

}