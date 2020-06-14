package UserUC

import (
	"github.com/mortawe/tech-db-forum/internal/models"
	"github.com/mortawe/tech-db-forum/internal/user"
)

type UserUC struct {
	repo user.IUserRepo
}

func NewUserUC(repo user.IUserRepo) *UserUC {
	return &UserUC{repo: repo}
}

func (uUC *UserUC) Insert(user *models.User) error {
	return uUC.repo.Insert(user)
}

func (uUC *UserUC) Update(user *models.User) error {
	return uUC.repo.Update(user)
}

func (uUC *UserUC) SelectByNickname(nickname string) (models.User, error) {
	return uUC.repo.SelectByNickname(nickname)
}

func (uUC *UserUC) SelectByEmail(email string) (models.User, error) {
	return uUC.repo.SelectByEmail(email)
}