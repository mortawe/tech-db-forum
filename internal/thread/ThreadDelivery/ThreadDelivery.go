package ThreadDelivery

import (
	"encoding/json"
	"github.com/fasthttp/router"
	"github.com/mortawe/tech-db-forum/internal/forum"
	"github.com/mortawe/tech-db-forum/internal/models"
	"github.com/mortawe/tech-db-forum/internal/thread"
	"github.com/mortawe/tech-db-forum/internal/user"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
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
	r.POST("/api/forum/{slug}/create", m.CreateThread)
	r.GET("/api/forum/{slug}/threads", m.GetThreadsByForum)
	r.GET("/api/thread/{slugOrID}/details", m.Details)
	r.POST("/api/thread/{slugOrID}/details", m.Update)
	r.POST("/api/thread/{slugOrID}/vote", m.Vote)
}

func (m *ThreadManager) CreateThread(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	logrus.Println("thread create")
	thread := &models.Thread{}

	if err := json.Unmarshal(ctx.PostBody(), thread); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(`{"message": "` + "unmarshal not ok" + `" }`))
		return
	}
	threadInDB, err := m.tUC.SelectThreadBySlug(thread.Slug)
	if thread.Slug != "" && err == nil {
		ctx.SetStatusCode(409)
		resp, _ := json.Marshal(threadInDB)
		ctx.Write(resp)
		return
	}
	thread.Forum = ctx.UserValue("slug").(string)
	forum, err := m.fUC.SelectBySlug(thread.Forum)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "` + "forum does not exists" + `" }`))
		return
	}
	thread.Forum = forum.Slug
	user, err := m.uUC.SelectByNickname(thread.Author)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "` + "user does not exists : " + err.Error() + " nickname : " + thread.Author + `" }`))
		return
	}
	thread.Author = user.Nickname

	err = m.tUC.InsertThread(thread)
	if err != nil {
		ctx.SetStatusCode(409)
		ctx.Write([]byte(`{"message": "` + err.Error() + `" }`))
		return
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)
	resp, _ := json.Marshal(thread)
	ctx.Write(resp)
	log.Print("success")
}

func (m *ThreadManager) GetThreadsByForum(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	logrus.Println("thread create")
	slug := ctx.UserValue("slug").(string)
	_, err := m.fUC.SelectBySlug(slug)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "` + "forum does not exists : " + err.Error() + `" }`))
		return
	}
	params := &models.GetThreadsParams{}
	params.Since = string(ctx.FormValue("since"))
	if string(ctx.FormValue("desc")) == "true" {
		params.Desc = true
	} else {
		params.Desc = false
	}
	params.Limit, _ = strconv.Atoi(string(ctx.FormValue("limit")))

	threads, err := m.tUC.SelectThreadsByForum(slug, params.Limit, params.Since, params.Desc)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(`{"message": "` + "select not ok : " + err.Error() + `" }`))
		return
	}
	resp, _ := json.Marshal(threads)
	ctx.SetStatusCode(200)
	ctx.Write(resp)
	log.Println("success")
}

func (m *ThreadManager) Details(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	logrus.Println("thread details")
	slug := ctx.UserValue("slugOrID").(string)
	thread, err := m.tUC.SelectBySlugOrID(slug)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "` + "thread does not exists : " + `" }`))
		return
	}
	thread.Votes, err = m.tUC.GetVoteCount(thread.ID)

	resp, _ := json.Marshal(thread)
	ctx.Write(resp)
	ctx.SetStatusCode(200)
}

func (m *ThreadManager) Update(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	logrus.Println("thread update")
	slug := ctx.UserValue("slugOrID").(string)
	threadInDB, err := m.tUC.SelectBySlugOrID(slug)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "` + "thread does not exists : " + `" }`))
		return
	}
	thread := &models.Thread{}
	if err := json.Unmarshal(ctx.PostBody(), thread); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(`{"message": "` + err.Error() + `" }`))
		return
	}
	thread.Slug = threadInDB.Slug
	if thread.Forum == "" {
		thread.Forum = threadInDB.Forum
	}
	if thread.Author == "" {
		thread.Author = threadInDB.Author
	}
	if thread.Title == "" {
		thread.Title = threadInDB.Title
	}
	if thread.Message == "" {
		thread.Message = threadInDB.Message
	}
	thread.ID = threadInDB.ID
	err = m.tUC.Update(thread)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(`{"message": "` + err.Error() + `" }`))
		return
	}
	resp, _ := json.Marshal(thread)
	ctx.Write(resp)
	ctx.SetStatusCode(200)
}



func (m *ThreadManager) Vote(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	slug := ctx.UserValue("slugOrID").(string)

	vote := &models.Vote{}
	if err := json.Unmarshal(ctx.PostBody(), vote); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(`{"message": "` + err.Error() + `" }`))
		return
	}
	_, err := m.uUC.SelectByNickname(vote.Nickname)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "` + "thread does not exists : " + `" }`))
		return
	}
	threadInDB, err := m.tUC.SelectBySlugOrID(slug)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "` + "thread does not exists : " + `" }`))
		return
	}
	if err := m.tUC.UpdateVoice(vote, threadInDB.ID); err != nil {
		er := m.tUC.InserteVoice(vote, threadInDB)
		if er != nil {
			ctx.Write([]byte(er.Error()))
			return
		}
	}
	threadInDB.Votes, err = m.tUC.GetVoteCount(threadInDB.ID)
	if err != nil {
		ctx.Write([]byte(err.Error()))
		return
	}
	resp , _ := json.Marshal(threadInDB)
	ctx.SetStatusCode(200)
	ctx.Write(resp)
}
