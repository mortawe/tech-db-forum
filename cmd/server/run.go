package main

import (
	"fmt"
	"github.com/fasthttp/router"
	"github.com/jackc/pgx"
	"github.com/mortawe/tech-db-forum/internal/forum/ForumDelivery"
	"github.com/mortawe/tech-db-forum/internal/forum/ForumRepo"
	"github.com/mortawe/tech-db-forum/internal/forum/ForumUC"
	"github.com/mortawe/tech-db-forum/internal/post/PostDelivery"
	"github.com/mortawe/tech-db-forum/internal/post/PostRepo"
	"github.com/mortawe/tech-db-forum/internal/post/PostUC"
	"github.com/mortawe/tech-db-forum/internal/service/ServiceDelivery"
	"github.com/mortawe/tech-db-forum/internal/thread/ThreadDelivery"
	"github.com/mortawe/tech-db-forum/internal/thread/ThreadRepo"
	"github.com/mortawe/tech-db-forum/internal/thread/ThreadUC"
	"github.com/mortawe/tech-db-forum/internal/user/UserDelivery"
	"github.com/mortawe/tech-db-forum/internal/user/UserRepo"
	"github.com/mortawe/tech-db-forum/internal/user/UserUC"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)


func main() {
	config := pgx.ConnPoolConfig{
		ConnConfig:     pgx.ConnConfig{
			Host:                 "localhost",
			Port:                 5432,
			Database:             "docker",
			User:                 "docker",
			Password:             "docker",
			TLSConfig:            nil,
			UseFallbackTLS:       false,
			FallbackTLSConfig:    nil,
			Logger:               nil,
			LogLevel:             0,
			Dial:                 nil,
			RuntimeParams:        nil,
			OnNotice:             nil,
			CustomConnInfo:       nil,
			CustomCancel:         nil,
			PreferSimpleProtocol: false,
			TargetSessionAttrs:   "",
		},
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}
	connPool, err := pgx.NewConnPool(config)
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Println("connected to db...")
	fR := ForumRepo.NewForumRepo(connPool)
	uR := UserRepo.NewUserRepo(connPool)
	tR := ThreadRepo.NewThreadRepo(connPool)
	pR := PostRepo.NewThreadRepo(connPool)

	fUC := ForumUC.NewForumUC(fR)
	uUC := UserUC.NewUserUC(uR)
	tUC := ThreadUC.NewForumUC(tR)
	pUC := PostUC.NewPostUC(pR)

	fM := ForumDelivery.NewForumManager(fUC, uUC)
	uM := UserDelivery.NewUserManager(uUC, fUC)
	tM := ThreadDelivery.NewThreadManager(fUC, uUC, tUC)
	pM := PostDelivery.NewForumManager(fUC, pUC, tUC, uUC)
	sM := ServiceDelivery.NewServiceManager(connPool)
	router := InitRoutes()
	fM.InitRoutes(router)
	uM.InitRoutes(router)
	tM.InitRoutes(router)
	pM.InitRoutes(router)
	sM.InitRouters(router)

	defer connPool.Close()
	fmt.Println("server started...")
	fasthttp.ListenAndServe(":5000", router.Handler)

}

func InitRoutes() *router.Router{
	r := router.New()

	return r
}
