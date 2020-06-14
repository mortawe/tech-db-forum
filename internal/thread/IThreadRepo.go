package thread

import "github.com/mortawe/tech-db-forum/internal/models"

type IThreadRepo interface {
	SelectThreadByTitle(title string) (*models.Thread, error)
	InsertThread(thread *models.Thread) error
	SelectThreadsByForum(slug string, limit int, since string, desc bool) ([]models.Thread, error)
	SelectThreadByID(id int) (*models.Thread, error)
	SelectThreadBySlug(slug string) (*models.Thread, error)
	Update(thread *models.Thread) error
	UpdateVoice(voice *models.Vote, thread int) (int, error)
	InsertVoice(voice *models.Vote, thread int) (int, error)
	GetVotes(thread int) (int, error)
}
