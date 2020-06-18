package thread

import "github.com/mortawe/tech-db-forum/internal/models"

type IThreadRepo interface {
	InsertThread(thread *models.Thread) error
	SelectThreadsByForum(slug string, limit int, since string, desc bool) ([]models.Thread, error)
	SelectThreadByID(id int) (*models.Thread, error)
	SelectThreadBySlug(slug string) (*models.Thread, error)
	Update(thread *models.Thread) error
	VoteBySlug(vote models.Vote, slug string) (models.Thread, error)
	VoteByID(vote models.Vote, id int) (models.Thread, error)
	UpdateBySlugOrID(s string, thread *models.Thread) error
	GetIDForumBySlugOrID(s string) (int, string, error)
	SelectForumByThreadID(id int) (string, error)
}
