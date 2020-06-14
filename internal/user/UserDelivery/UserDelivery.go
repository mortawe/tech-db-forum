package UserDelivery

import (
	"encoding/json"
	"github.com/fasthttp/router"
	"github.com/mortawe/tech-db-forum/internal/forum"
	"github.com/mortawe/tech-db-forum/internal/models"
	"github.com/mortawe/tech-db-forum/internal/user"
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
	ctx.SetContentType("application/json")
	if err := json.Unmarshal(ctx.PostBody(), user); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("user create not ok when unmarshal : " + err.Error()))
		return
	}
	user.Nickname = ctx.UserValue("nickname").(string)
	if user.Nickname == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("user create not ok when get nickname : "+ user.Nickname))
		return
	}
	arr := []models.User{}
	userInDB, err := m.uUC.SelectByNickname(user.Nickname)
	if err == nil {
		ctx.SetStatusCode(409)
		arr = append(arr, userInDB)
	}
	userInDBEmail, err := m.uUC.SelectByEmail(user.Email)
	if err == nil && userInDBEmail.Nickname != userInDB.Nickname {
		ctx.SetStatusCode(409)
		arr = append(arr, userInDBEmail)
	}

	if len(arr) > 0 {
		resp, _ := json.Marshal(arr)
		ctx.Write(resp)
		return
	}
	if err := m.uUC.Insert(user); err != nil {
		ctx.SetStatusCode(409)

		arr := []models.User{userInDB}
		resp, _ := json.Marshal(arr)
		ctx.Write(resp)
		return
	}
	resp, _ := json.Marshal(user)
	ctx.Write(resp)
	ctx.SetStatusCode(201)

}

func (m *UserManager) UpdateProfile(ctx *fasthttp.RequestCtx) {
	user := &models.User{}
	ctx.SetContentType("application/json")
	if err := json.Unmarshal(ctx.PostBody(), user); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(`{"message": "Can't unmarshal'"}`))
		return
	}
	user.Nickname = ctx.UserValue("nickname").(string)
	if user.Nickname == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(`{"message": "user nickname is empty"}`))
		return
	}
	userInDB, err := m.uUC.SelectByNickname(user.Nickname)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "can't select by nickname'"}`))
		return
	}
	if userInDB.Nickname != user.Nickname {
		ctx.SetStatusCode(409)
		ctx.Write([]byte(`{"message": "` + userInDB.Nickname + " " + user.Nickname + `"}`) )
		return
	}
	if user.Email == "" {
		user.Email = userInDB.Email
	}
	if user.About == "" {
		user.About = userInDB.About
	}
	if user.Fullname == "" {
		user.Fullname = userInDB.Fullname
	}
	if _, err := m.uUC.SelectByEmail(user.Email); err == nil && userInDB.Email != user.Email {
		ctx.SetStatusCode(409)
		ctx.Write([]byte(`{"message": "` + userInDB.Email + " " + user.Email + `"}`))
		return
	}

	if err := m.uUC.Update(user); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(`{"message": "Can't update : ` + err.Error() + `"}`))
		return
	}
	ctx.SetStatusCode(200)
	resp, _ := json.Marshal(user)
	ctx.Write(resp)
}

func (m *UserManager) GetProfile(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	user := &models.User{}

	user.Nickname = ctx.UserValue("nickname").(string)
	if user.Nickname == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("not ok when get nickname : " + user.Nickname))
		return
	}
	userInDB, err := m.uUC.SelectByNickname(user.Nickname)
	if err != nil {
		ctx.SetStatusCode(404)
		ctx.Write([]byte(`{"message": "Can't'"}`))
		return
	}
	resp, _ := json.Marshal(userInDB)
	ctx.Write(resp)
	ctx.SetStatusCode(200)
}
