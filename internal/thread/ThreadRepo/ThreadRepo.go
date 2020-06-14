package ThreadRepo

import (
	"github.com/jackc/pgx"
	"github.com/mortawe/tech-db-forum/internal/models"
)

type ThreadRepo struct {
	db *pgx.ConnPool
}

func NewThreadRepo(db *pgx.ConnPool) *ThreadRepo {
	return &ThreadRepo{
		db: db,
	}
}

func (r *ThreadRepo) InsertThread(thread *models.Thread) error {
	err := r.db.QueryRow("INSERT INTO threads (author, created, forum, message, slug, title) " +
		"VALUES ($1, $2, $3, $4, $5, $6) " +
		"RETURNING id, votes", &thread.Author, &thread.Created, &thread.Forum, &thread.Message, &thread.Slug,
		&thread.Title).Scan(&thread.ID, &thread.Votes)
	return err
}

func (r *ThreadRepo) SelectThreadByTitle(title string) (*models.Thread, error) {
	thread := models.Thread{}
	err := r.db.QueryRow("SELECT * FROM threads WHERE title = $1", title).Scan(&thread.Author, &thread.Created,
		&thread.Forum, &thread.ID, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	if err != nil {
		return nil, err
	}
	return &thread, nil
}

func (r *ThreadRepo) SelectThreadByID(id int) (*models.Thread, error) {
	thread := &models.Thread{}
	err := r.db.QueryRow("SELECT * FROM threads WHERE id = $1", id).Scan(&thread.Author, &thread.Created,
		&thread.Forum, &thread.ID, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (r *ThreadRepo) SelectThreadsByForum(slug string, limit int, since string, desc bool) ([]models.Thread, error) {
	threads := []models.Thread{}

	query := "SELECT * FROM threads WHERE forum = $1 "
	rows := &pgx.Rows{}
	var err error
	if limit > 0  && since != "" {
		if desc {
			query += "AND created <= $2 ORDER BY created DESC LIMIT $3"
		} else {
			query += "AND created >= $2 ORDER BY created ASC LIMIT $3"
		}
		rows, err = r.db.Query(query, &slug, &since, &limit)
	} else {
		if limit > 0 {
			if desc {
				query += "ORDER BY created DESC LIMIT $2"
			} else {
				query += "ORDER BY created ASC LIMIT $2"
			}
			rows, err = r.db.Query(query, &slug, &limit)
		}
		if since != "" {
			if desc {
				query += "AND created <= $2 ORDER BY created DESC"
			} else {
				query += "AND created >= $2 ORDER BY created ASC"
			}
			rows, err = r.db.Query(query, &slug, &since)
		}
	}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		thread := models.Thread{}

		err := rows.Scan(&thread.Author, &thread.Created,&thread.Forum,&thread.ID, &thread.Message, &thread.Slug,
			&thread.Title, &thread.Votes);
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	return threads, nil
}

func (r *ThreadRepo) SelectThreadBySlug(slug string) (*models.Thread, error) {
	thread := models.Thread{}
	err := r.db.QueryRow("SELECT * FROM threads WHERE slug = $1", slug).Scan(&thread.Author, &thread.Created,
		&thread.Forum, &thread.ID, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	if err != nil {
		return nil, err
	}
	return &thread, nil
}


func (r *ThreadRepo) Update(thread *models.Thread) error {
	err := r.db.QueryRow("UPDATE threads "+
		"SET author = $1, "+
		"forum = $2, "+
		"message = $3, " +
		"slug = $4, " +
		"title = $5  " +
		"WHERE id = $6 " +
		"RETURNING threads.* ", &thread.Author, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title, &thread.ID).Scan(
			&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message, &thread.Slug,
		&thread.Title, &thread.Votes)
	return err
}

func (r *ThreadRepo) VoteBySlug(vote models.Vote, slug string) (models.Thread, error) {
	thread, err := r.SelectThreadBySlug(slug)
	if err != nil {
		return models.Thread{}, err
	}
	_, err = r.db.Exec(`INSERT INTO votes (threadid, nickname, vote)
			VALUES ($1, $2, $3)
			ON CONFLICT (threadid, nickname) DO UPDATE SET vote = $3`,
		thread.ID,
		vote.Nickname,
		vote.Voice,
	)
	if err != nil {
		return models.Thread{}, err
	}
	err = r.db.QueryRow(
		`SELECT votes FROM threads WHERE id = $1`,
		thread.ID).Scan(&thread.Votes)
	if err != nil {
		return models.Thread{}, err
	}
	return *thread, nil
}

func (r *ThreadRepo) VoteByID(vote models.Vote, id int) (models.Thread, error) {
	_, err := r.db.Exec(`
			INSERT INTO votes (threadid, nickname, vote)
			VALUES ($1, $2, $3)
			ON CONFLICT (threadid, nickname) DO UPDATE SET vote = $3`,
		id,
		vote.Nickname,
		vote.Voice,
	)
	if err != nil {
		return models.Thread{}, err
	}
	thread := &models.Thread{}
	if thread, err = r.SelectThreadByID(id); err != nil {
		return models.Thread{}, err
	}
	return *thread, nil
}
