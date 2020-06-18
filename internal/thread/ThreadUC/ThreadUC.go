package ThreadUC

import (
	"github.com/mortawe/tech-db-forum/internal/models"
	"github.com/mortawe/tech-db-forum/internal/thread"
	"strconv"
)

type ThreadUC struct {
	repo thread.IThreadRepo
}

func NewForumUC(repo thread.IThreadRepo) *ThreadUC {
	return &ThreadUC{repo: repo}
}

func (uc *ThreadUC) InsertThread(thread *models.Thread) error {
	return uc.repo.InsertThread(thread)
}

func (uc *ThreadUC) SelectThreadsByForum(slug string, limit int, since string, desc bool) ([]models.Thread, error) {
	return uc.repo.SelectThreadsByForum(slug, limit, since, desc)
}

func (uc *ThreadUC) SelectBySlugOrID(slugOrID string) (*models.Thread, error) {
	value, err := strconv.Atoi(slugOrID)
	if err == nil {
		return uc.repo.SelectThreadByID(value)
	} else {
		return uc.repo.SelectThreadBySlug(slugOrID)
	}
}
func (uc *ThreadUC) SelectByID(id int) (*models.Thread, error) {
	return uc.repo.SelectThreadByID(id)
}

func (uc *ThreadUC) SelectThreadBySlug(slug string) (*models.Thread, error) {
	return uc.repo.SelectThreadBySlug(slug)
}
func (uc *ThreadUC) Update(thread *models.Thread) error {
	return uc.repo.Update(thread)
}

func (uc *ThreadUC) Vote(vote models.Vote, slug string) (models.Thread, error) {
	if id, err := strconv.Atoi(slug); err != nil {
		return uc.repo.VoteBySlug(vote, slug)
	} else {
		return uc.repo.VoteByID(vote, id)
	}
}

func (uc *ThreadUC) UpdateBySlugOrID(s string, thread *models.Thread) error {
	return uc.repo.UpdateBySlugOrID(s, thread)
}

func (uc *ThreadUC) GetIDForumBySlugOrID(s string) (int, string, error) {
	value, err := strconv.Atoi(s)
	if err == nil {
		forum, err := uc.repo.SelectForumByThreadID(value)
		return value, forum, err
	} else {
		return uc.repo.GetIDForumBySlugOrID(s)
	}
}