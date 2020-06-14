package forum

import "github.com/mortawe/tech-db-forum/internal/models"

type IForumUC interface {
	Create(forum *models.Forum) error
	SelectBySlug(slug string) (*models.Forum, error)
	GetUsersByForum(slug string, desc bool, since string, limit int) ([]models.User, error)
}
