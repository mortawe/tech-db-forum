package PostRepo

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/mortawe/tech-db-forum/internal/models"
)
type PostRepo struct {
	db *pgx.ConnPool
}

func NewThreadRepo(db *pgx.ConnPool) *PostRepo {
	return &PostRepo{
		db: db,
	}
}

func (r *PostRepo) SelectThreadByPostID(id int) (int, error) {
	tID := 0
	err := r.db.QueryRow("SELECT thread FROM posts WHERE id = $1", &id).Scan(&tID)
	return tID, err
}

func (r *PostRepo) InsertPosts(posts []*models.Post) error {
	query := "INSERT INTO posts (author, forum, message, parent, thread) values "
	if len(posts) == 0 {
		return nil
	}
	for i, p := range posts {
		if i != 0 {
			query += ", "
		}
		query += fmt.Sprintf("('%s', '%s', '%s', %d, %d) ", p.Author, p.Forum, p.Message,
			p.Parent,p.Thread)
	}

	query += "RETURNING id, created"
	rows, err := r.db.Query(query)
	if err != nil {
		return err
	}
	for idx := 0; rows.Next(); idx++ {
		if err := rows.Scan(&posts[idx].ID, &posts[idx].Created); err != nil {
			return err
		}
	}
	return nil
}

func (r *PostRepo) SelectPostByID(id int) (*models.Post, error) {
	post := &models.Post{}
	err := r.db.QueryRow("SELECT author, created, forum, id, edited, message, parent, thread FROM posts WHERE id = " +
		"$1", &id).Scan(&post.Author, &post.Created, &post.Forum,
		&post.ID, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
	return post, err
}

func (r *PostRepo) Update(post *models.Post) error {
	err := r.db.QueryRow("UPDATE posts SET message = $1, edited = true WHERE id = $2 RETURNING author, created, forum, id, edited," +
		" message, parent, thread", post.Message, post.ID).Scan(&post.Author, &post.Created, &post.Forum,
			&post.ID, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
	return err
}

func (r *PostRepo) GetPosts(threadID int, desc bool, since string, limit int, sort string) ([]models.Post, error){
	posts := []models.Post{}
	query := ""

	var err error
	rows := &pgx.Rows{}
	if since != "" {
		switch sort {
		case "tree":
			query = "SELECT posts.id, posts.author, posts.forum, posts.thread, " +
				"posts.message, posts.parent, posts.edited, posts.created " +
				"FROM posts %s posts.thread = $1 ORDER BY posts.path[1] %s, posts.path %s LIMIT $3"
			if desc {
				query = fmt.Sprintf(query, "JOIN posts P ON P.id = $2 WHERE posts.path < p.path AND",
					"DESC",
					"DESC")
			} else {
				query = fmt.Sprintf(query, "JOIN posts P ON P.id = $2 WHERE posts.path > p.path AND",
					"ASC",
					"ASC")
			}
		case "parent_tree":
			query =  "SELECT p.id, p.author, p.forum, p.thread, p.message, p.parent, p.edited, p.created " +
				"FROM posts as p WHERE p.thread = $1 AND " +
				"p.path::integer[] && (SELECT ARRAY (select p.id from posts as p WHERE p.thread = $1 AND p.parent = 0 %s %s %s"
			if desc {
				query = fmt.Sprintf(query, " AND p.path < (SELECT p.path[1:1] FROM posts as p WHERE p.id = $2) ",
					"ORDER BY p.path[1] DESC, p.path LIMIT $3)) ",
					"ORDER BY p.path[1] DESC, p.path ")
			} else {
				query = fmt.Sprintf(query, " AND p.path > (SELECT p.path[1:1] FROM posts as p WHERE p.id = $2) ",
					"ORDER BY p.path[1] ASC, p.path LIMIT $3)) ",
					"ORDER BY p.path[1] ASC, p.path ")
			}
		default:
			query = "SELECT id, author, forum, thread, message, parent, edited, created " +
				"FROM posts WHERE thread = $1 AND id %s $2 ORDER BY id %s LIMIT $3"
			if desc {
				query = fmt.Sprintf(query, "<", "DESC")
			} else {
				query = fmt.Sprintf(query, ">", "ASC")
			}
		}
		rows, err = r.db.Query(query, threadID, since, limit)
	} else {
		switch sort {
		case "tree":
			if desc {
				query = fmt.Sprintf("SELECT posts.id, posts.author, posts.forum, posts.thread, " +
					"posts.message, posts.parent, posts.edited, posts.created " +
					"FROM posts WHERE posts.thread = $1 ORDER BY posts.path[1] DESC, posts.path DESC LIMIT $2")
			} else {
				query = fmt.Sprintf("SELECT posts.id, posts.author, posts.forum, posts.thread, " +
					"posts.message, posts.parent, posts.edited, posts.created " +
					"FROM posts WHERE posts.thread = $1 ORDER BY posts.path[1] ASC, posts.path ASC LIMIT $2")
			}
		case "parent_tree":
			if desc {
				query = "SELECT p.id, p.author, p.forum, p.thread, p.message, p.parent, p.edited, p.created " +
					"FROM posts as p WHERE p.thread = $1 AND " +
					"p.path::integer[] && (SELECT ARRAY (select p.id from posts as p WHERE p.thread = $1 AND p.parent = 0" +
					"ORDER BY p.path[1] DESC, p.path LIMIT $2)) " +
					"ORDER BY p.path[1] DESC, p.path"
			} else {
				query ="SELECT p.id, p.author, p.forum, p.thread, p.message, p.parent, p.edited, p.created " +
					"FROM posts as p WHERE p.thread = $1 AND " +
					"p.path::integer[] && (SELECT ARRAY (select p.id from posts as p WHERE p.thread = $1 AND p.parent = 0 " +
					"ORDER BY p.path[1] ASC, p.path LIMIT $2)) ORDER BY p.path[1] ASC, p.path"
			}
		default:
			if desc {
				query = "SELECT id, author, forum, thread, message, parent, edited, created " +
					"FROM posts WHERE thread = $1  ORDER BY id DESC LIMIT $2"
			} else {
				query = "SELECT id, author, forum, thread, message, parent, edited, created " +
					"FROM posts WHERE thread = $1 ORDER BY id ASC LIMIT $2"
			}
		}
		rows, err = r.db.Query(query, threadID, limit)
	}

	if err != nil {
		return posts, err
	}

	for rows.Next() {
		p := &models.Post{}
		err := rows.Scan(&p.ID, &p.Author, &p.Forum, &p.Thread, &p.Message, &p.Parent, &p.IsEdited, &p.Created)
		if err != nil {
			return posts, err
		}
		posts = append(posts, *p)
	}
	return posts, nil
}