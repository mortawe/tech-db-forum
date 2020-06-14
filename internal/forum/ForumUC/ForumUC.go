package ForumUC

import (
	"github.com/mortawe/tech-db-forum/internal/forum"
	"github.com/mortawe/tech-db-forum/internal/models"
)

type ForumUC struct {
	repo forum.IForumRepo
}

func NewForumUC(repo forum.IForumRepo) *ForumUC {
	return &ForumUC{repo: repo}
}

func (uc *ForumUC) Create(forum *models.Forum) error {
	return uc.repo.Insert(forum)
}

func (uc *ForumUC) SelectBySlug(slug string) (*models.Forum, error) {
	return uc.repo.SelectBySlug(slug)
}

func (uc *ForumUC) GetUsersByForum(slug string, desc bool, since string, limit int) ([]models.User, error) {
	return uc.repo.GetUsersByForum(slug, desc, since, limit)
}

