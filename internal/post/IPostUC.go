package post

import (
	"github.com/mortawe/tech-db-forum/internal/models"
	"time"
)

type IPostUC interface {
	InsertPost(posts []*models.Post, forum string, threadID int, created time.Time) error
	SelectPostByID(id int) (*models.Post, error)
	Update(post *models.Post) error
	GetPosts(threadID int, desc bool, since string, limit int, sort string) ([]models.Post, error)



}
