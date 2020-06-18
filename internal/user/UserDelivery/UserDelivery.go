package UserDelivery

import (
	"encoding/json"
	"errors"
	"github.com/fasthttp/router"
	"github.com/mortawe/tech-db-forum/internal/forum"
	"github.com/mortawe/tech-db-forum/internal/models"
	"github.com/mortawe/tech-db-forum/internal/user"
	"github.com/mortawe/tech-db-forum/internal/utils"
	"github.com/valyala/fasthttp"
)

type UserManager struct {
	uUC user.IUserUC
	fUC forum.IForumUC
}

func NewUserManager(u user.IUserRepo, f forum.IForumUC) *UserManager {
	return &UserManager{uUC: u, fUC: f}
}

func (m *UserManager) InitRoutes(r *router.Router) {
	r.POST("/api/user/{nickname}/create", m.CreateUser)
	r.POST("/api/user/{nickname}/profile", m.UpdateProfile)
	r.GET("/api/user/{nickname}/profile", m.GetProfile)
}

func (m *UserManager) CreateUser(ctx *fasthttp.RequestCtx) {
	user := &models.User{}
	if err := user.UnmarshalJSON(ctx.PostBody()); err != nil {
		utils.Send(400, ctx, utils.MustMarshalError(err))
		return
	}
	user.Nickname = ctx.UserValue("nickname").(string)
	if user.Nickname == "" {
		utils.Send(400, ctx, utils.MustMarshalError(errors.New("no nickname")))
		return
	}
	err := m.uUC.Insert(user)
	switch err {
	case models.ErrConflict:
		usersAlreadyExist, err := m.uUC.SelectByEmailOrNickname(user.Nickname, user.Email)
		if err != nil {
			utils.Send(400, ctx, utils.MustMarshalError(err))
			return
		}
		resp, _ := json.Marshal(usersAlreadyExist)
		utils.Send(409, ctx, resp)
	case nil:
		resp, _ := user.MarshalJSON()
		utils.Send(201, ctx, resp)
	default:
		utils.Send(500, ctx, utils.MustMarshalError(err))
		return
	}


}

func (m *UserManager) UpdateProfile(ctx *fasthttp.RequestCtx) {
	user := &models.User{}
	if err := user.UnmarshalJSON(ctx.PostBody()); err != nil {
		utils.Send(400, ctx, utils.MustMarshalError(err))
		return
	}
	user.Nickname = ctx.UserValue("nickname").(string)
	if user.Nickname == "" {
		utils.Send(400, ctx, utils.MustMarshalError(errors.New("no nickname")))
		return
	}
	err := m.uUC.Update(user)
	switch err {
	case models.ErrConflict:
		utils.Send(409, ctx, utils.MustMarshalError(err))
	case models.ErrNotExists:
		utils.Send(404, ctx, utils.MustMarshalError(err))
	case nil:
		resp, _ := user.MarshalJSON()
		utils.Send(200, ctx, resp)
	default:
		utils.Send(500, ctx, utils.MustMarshalError(err))
	}
}

func (m *UserManager) GetProfile(ctx *fasthttp.RequestCtx) {
	user := &models.User{}

	user.Nickname = ctx.UserValue("nickname").(string)
	if user.Nickname == "" {
		utils.Send(400, ctx, []byte(`Can't get nickname`))
		return
	}
	userInDB, err := m.uUC.SelectByNickname(user.Nickname)
	switch err {
	case models.ErrNotExists :
		utils.Send(404, ctx, utils.MustMarshalError(err))
		return
	case nil:
		resp, _ := userInDB.MarshalJSON()
		utils.Send(200, ctx, resp)
	default:
		utils.Send(500, ctx, utils.MustMarshalError(err))
	}
}
