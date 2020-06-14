package PostUC

import (
	"errors"
	"github.com/mortawe/tech-db-forum/internal/models"
	"github.com/mortawe/tech-db-forum/internal/post"
	"time"
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

func (uc *PostUC) InsertPost(posts []*models.Post, forum string, threadID int, created time.Time) error {
	for _, post := range posts {
		if post.Parent != 0 {
			parentPostThreadID, err := uc.repo.SelectThreadByPostID(post.Parent)
			if err != nil {
				return ParentErr
			}
			if parentPostThreadID != threadID {
				return ParentErr
			}
		}
		post.Forum = forum
		post.Thread = threadID
		post.Created = created
		post.IsEdited = false
	}
	err := uc.repo.InsertPosts(posts)
	return err
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