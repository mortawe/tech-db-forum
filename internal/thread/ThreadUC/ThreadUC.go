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

func (uc *ThreadUC) InsertVoice(voice *models.Vote, thread int) (int, error){
	return uc.repo.InsertVoice(voice, thread)
}
func (uc *ThreadUC) UpdateVoice(voice *models.Vote, thread int) (int, error) {
	return uc.repo.UpdateVoice(voice, thread)
}

func (uc *ThreadUC) GetVotes(thread int) (int, error) {
	return uc.repo.GetVotes(thread)
}