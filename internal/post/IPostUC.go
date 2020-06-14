package post

import (
	"github.com/mortawe/tech-db-forum/internal/models"
)

type IPostUC interface {
	InsertPost(posts []*models.Post, forum string, threadID int) error
	SelectPostByID(id int) (*models.Post, error)
	Update(post *models.Post) error
	GetPosts(threadID int, desc bool, since string, limit int, sort string) ([]models.Post, error)



}
