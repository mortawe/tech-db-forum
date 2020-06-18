package user

import "github.com/mortawe/tech-db-forum/internal/models"

type IUserRepo interface {
	Insert(user *models.User) error
	Update(user *models.User) error
	SelectByNickname(nickname string) (models.User, error)
	SelectNicknameWithCase(nickname string) (string, error)
	SelectByEmailOrNickname(nickname string, email string) (models.Users, error)
}