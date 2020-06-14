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
func (uc *ThreadUC) SelectThreadByTitle(title string) (*models.Thread, error) {
	return uc.repo.SelectThreadByTitle(title)
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

func (uc *ThreadUC) InserteVoice(voice *models.Vote, thread *models.Thread) error {
	return uc.repo.InserteVoice(voice, thread)
}
func (uc *ThreadUC) UpdateVoice(voice *models.Vote, thread int) error{
	return uc.repo.UpdateVoice(voice, thread)
}
func (uc *ThreadUC)  SelecteVoice(nickname string, thread int) (*models.Vote, error) {
	return uc.repo.SelecteVoice(nickname, thread)
}

func (uc *ThreadUC) GetVoteCount(id int) (int, error) {
	return uc.repo.GetVoteCount(id)
}