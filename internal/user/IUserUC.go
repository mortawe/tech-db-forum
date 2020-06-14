package user

import "github.com/mortawe/tech-db-forum/internal/models"

type IUserUC interface {
	Insert(user *models.User) error
	Update(user *models.User) error
	SelectByNickname(nickname string) (models.User, error)
	SelectByEmail(email string) (models.User, error)
}