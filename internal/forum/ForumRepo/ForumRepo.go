package ForumRepo

import (
	"github.com/jackc/pgx"
	"github.com/mortawe/tech-db-forum/internal/models"
	"github.com/sirupsen/logrus"
)

type ForumRepo struct {
	db *pgx.ConnPool
}

func NewForumRepo(db *pgx.ConnPool) *ForumRepo {
	return &ForumRepo{
		db: db,
	}
}

func (r *ForumRepo) Insert(forum *models.Forum) error {
	err := r.db.QueryRow("INSERT INTO forums (slug, title, nickname) " +
		"VALUES ($1, $2, $3) " +
		"RETURNING forums.* " +
		"", &forum.Slug, &forum.Title, &forum.User).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Posts, &forum.Threads)

	return err
}

func (r *ForumRepo) SelectBySlug(slug string) (*models.Forum, error) {
	forum := &models.Forum{}
	err := r.db.QueryRow("SELECT forums.* " +
		"FROM forums " +
		"WHERE slug = $1 ", slug).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Posts, &forum.Threads)
	if err != nil {
		return nil, err
	}
	return forum, nil
}

func (r *ForumRepo) GetUsersByForum(slug string, desc bool, since string, limit int) ([]models.User, error){
	users := []models.User{}
	query := "SELECT DISTINCT users.* FROM user_forum JOIN users on users.nickname = user_forum.author WHERE forum = $1 "
	rows := &pgx.Rows{}
	var err error
	if limit > 0  && since != "" {
		if desc {
			query += "AND lower(users.nickname) < lower($2::text) ORDER BY users.nickname  DESC LIMIT $3"
		} else {
			query += "AND lower(users.nickname)  > lower($2::text) ORDER BY users.nickname  ASC LIMIT $3"
		}
		rows, err = r.db.Query(query, &slug, &since, &limit)
	} else {
		if limit > 0 {
			if desc {
				query += "ORDER BY users.nickname DESC LIMIT $2"
			} else {
				query += "ORDER BY users.nickname ASC LIMIT $2"
			}
			rows, err = r.db.Query(query, &slug, &limit)
		} else
		if since != "" {
			if desc {
				query += "AND lower(users.nickname) < lower($2::text) ORDER BY users.nickname DESC "
			} else {
				query += "AND lower(users.nickname) > lower($2::text) ORDER BY users.nickname ASC "
			}
			rows, err = r.db.Query(query, &slug, &since)
		} else {
			rows, err = r.db.Query(query, &slug)
		}
	}
	logrus.Println(query, " ", slug, " ", since, " ")
	if err != nil {
		logrus.Println("THE GREAT ERROR : ", err)
		return users, nil
	}
	for rows.Next() {
		user := models.User{}

		err := rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname);
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return  users, nil
}
