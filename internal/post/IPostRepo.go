package post

import "github.com/mortawe/tech-db-forum/internal/models"

type IPostRepo interface {
	SelectThreadByPostID(id int) (int, error)
	InsertPost( posts []*models.Post, forum string, id int) error
	SelectPostByID(id int) (*models.Post, error)
	Update(post *models.Post) error
	GetPosts(threadID int, desc bool, since string, limit int, sort string) ([]models.Post, error)
}
