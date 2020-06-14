package ForumDelivery

import (
	"encoding/json"
	"github.com/fasthttp/router"
	"github.com/mortawe/tech-db-forum/internal/forum"
	"github.com/mortawe/tech-db-forum/internal/models"
	"github.com/mortawe/tech-db-forum/internal/user"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"strconv"
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
	ctx.SetContentType("application/json")

	logrus.Println("forum create")
	forum := &models.Forum{}

	if err := json.Unmarshal(ctx.PostBody(), forum); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(`{"message": "` + "umarshal" + `"}`))
		return
	}
	user, _ := m.uUC.SelectByNickname(forum.User)
	forum.User = user.Nickname
	forumInBase, err := m.fUC.SelectBySlug(forum.Slug)
	if err == nil {
		ctx.SetStatusCode(409)
		resp, _ := json.Marshal(forumInBase)
		ctx.Write(resp)
		return
	}
	//log.Println(err)
	err = m.fUC.Create(forum)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "` + `"}`))
		return
	}
	resp, _ := json.Marshal(forum)
	ctx.Write(resp)
	ctx.SetStatusCode(fasthttp.StatusCreated)
}

func (m *ForumManager) Details(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slug := ctx.UserValue("slug").(string)
	resp := []byte("")

	forum, err := m.fUC.SelectBySlug(slug)

	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "` + "exit forum detailes bad req" + `"}`))
		return
	}

	resp, _ = json.Marshal(forum)
	ctx.SetStatusCode(200)
	ctx.Write(resp)
}

func (m *ForumManager) GetUsersByForum(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	params := &models.GetThreadsParams{}
	params.Since = string(ctx.FormValue("since"))
	if string(ctx.FormValue("desc")) == "true" {
		params.Desc = true
	} else {
		params.Desc = false
	}
	params.Limit, _ = strconv.Atoi(string(ctx.FormValue("limit")))

	slug := ctx.UserValue("slug").(string)
	thread, err := m.fUC.SelectBySlug(slug)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	users, err := m.fUC.GetUsersByForum(thread.Slug, params.Desc, params.Since, params.Limit)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	resp, _ := json.Marshal(users)
	ctx.SetStatusCode(200)
	ctx.Write(resp)
}