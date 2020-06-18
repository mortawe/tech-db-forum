package PostUC

import (
	"errors"
	"github.com/mortawe/tech-db-forum/internal/models"
	"github.com/mortawe/tech-db-forum/internal/post"
)

var (
	ParentErr = errors.New("parent")
)

type PostUC struct {
	repo post.IPostRepo
}

func NewPostUC(repo post.IPostRepo) *PostUC {
	return &PostUC{repo: repo}
}

func (uc *PostUC) InsertPost(posts []*models.Post, forum string, id int) error {
	return uc.repo.InsertPost(posts, forum, id)
}

func (uc *PostUC) Update(post *models.Post) error {
	return uc.repo.Update(post)
}

func (uc *PostUC) SelectPostByID(id int) (*models.Post, error) {
	return uc.repo.SelectPostByID(id)
}

func (uc *PostUC) GetPosts(threadID int, desc bool, since string, limit int, sort string) ([]models.Post, error) {
	return uc.repo.GetPosts(threadID, desc, since, limit, sort)
}
