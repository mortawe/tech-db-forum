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
func (r *ThreadRepo) GetVoteCount(id int) (int, error) {
	thread := models.Thread{}
	err := r.db.QueryRow("SELECT SUM(vote) FROM votes where threadid = $1", id).Scan(&thread.Votes)
	return thread.Votes, err
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

func (r *ThreadRepo) InserteVoice(voice *models.Vote, thread *models.Thread) error {
	_, err := r.db.Exec("INSERT INTO votes (threadid, nickname, vote) VALUES ($1, $2, $3)", thread.ID, voice.Nickname, voice.Voice)
	return err
}

func (r *ThreadRepo) UpdateVoice(voice *models.Vote, thread int) error {
	threaded := models.Thread{}
	err := r.db.QueryRow("UPDATE votes SET vote = $1 WHERE nickname = $2 AND threadid = $3 RETURNING vote", voice.Voice, voice.Nickname, thread).Scan(&threaded.Votes)
	return err
}

func (r *ThreadRepo) SelecteVoice(nickname string, thread int) (*models.Vote, error) {
	vote := &models.Vote{}
	err := r.db.QueryRow("SELECT  vote, nickname FROM votes " +
		"WHERE nickname = $1 AND threadid = $2", nickname, thread).Scan(&vote.Voice,
			&vote.Nickname)
	return vote, err
}