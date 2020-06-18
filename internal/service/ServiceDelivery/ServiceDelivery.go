package ServiceDelivery

import (
	"github.com/fasthttp/router"
	"github.com/jackc/pgx"
	"github.com/mortawe/tech-db-forum/internal/models"
	"github.com/valyala/fasthttp"
)

type ServiceManager struct {
	db *pgx.ConnPool
}

func NewServiceManager(db *pgx.ConnPool) *ServiceManager {
	return &ServiceManager{db: db}
}

func (m *ServiceManager) InitRouters(r *router.Router) {
	r.GET("/api/service/status", m.Status)
	r.POST("/api/service/clear", m.Clear)
}

func (m *ServiceManager) Status(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	status := &models.Status{}
	err := m.db.QueryRow("SELECT " +
		"(SELECT COUNT(*) FROM forums) as forums_status, " +
		"(SELECT COUNT(*) FROM threads) as threads_status, " +
		"(SELECT COUNT(*) FROM posts) as posts_status, " +
		"(SELECT COUNT(*) FROM users) as users_status").Scan(
			&status.Forums, &status.Threads, &status.Posts,
			&status.Users)
	if err != nil {

	}
	resp, _ := status.MarshalJSON()
	ctx.Write(resp)
	ctx.SetStatusCode(200)
}

func (m *ServiceManager) Clear(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	m.db.Exec("TRUNCATE  votes, posts, forum_users, threads, forums, users CASCADE")
	ctx.SetStatusCode(200)
}
