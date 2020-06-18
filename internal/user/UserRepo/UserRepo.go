package UserRepo

import (
	"github.com/jackc/pgx"
	"github.com/mortawe/tech-db-forum/internal/models"
)

type UserRepo struct {
	db *pgx.ConnPool
}

func NewUserRepo(db *pgx.ConnPool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Insert(user *models.User) error {
	_, err := r.db.Exec("INSERT INTO users "+
		"VALUES ($1, $2, $3, $4)",
		user.Nickname,
		user.Email,
		user.About,
		user.Fullname)
	switch err {
	case nil:
		return nil
	default:
		return models.ErrConflict
	}
}

func (r *UserRepo) Update(user *models.User) error {
	query := "UPDATE users SET   "
	if user.About != "" {
		query += ` about = '` + user.About + "' , "
	}
	if user.Fullname != "" {
		query += " fullname = '" + user.Fullname + "' , "
	}
	if user.Email != "" {
		query += " email = '" + user.Email + "' , "
	}
	query = query[:len(query)-2]
	query += "WHERE nickname = '" + user.Nickname + "' RETURNING about, email, fullname, nickname "
	if user.About == "" && user.Email == "" && user.Fullname == "" {
		query = "SELECT about, email, fullname, nickname FROM users WHERE nickname = '" + user.Nickname + "' "
	}
	err := r.db.QueryRow(query).Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return models.ErrNotExists
		default:
			return models.ErrConflict
		}
	}
	return nil
}

func (r *UserRepo) SelectByNickname(nickname string) (models.User, error) {
	user := models.User{}
	err := r.db.QueryRow("SELECT * FROM users "+
		"WHERE nickname = $1 ", nickname).Scan(&user.Nickname, &user.Email, &user.About, &user.Fullname)
	switch err {
	case pgx.ErrNoRows:
		return models.User{}, models.ErrNotExists
	default:
		return user, err
	}
}

func (r *UserRepo) SelectByEmailOrNickname(nickname string, email string) (models.Users, error) {
	users := []models.User{}
	rows, err := r.db.Query("SELECT * FROM users "+
		"WHERE nickname = $1 OR email = $2 LIMIT 2", nickname, email)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		user := models.User{}
		rows.Scan(&user.Nickname, &user.Email, &user.About, &user.Fullname)
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepo) SelectNicknameWithCase(nickname string) (string, error) {
	result := ""
	err := r.db.QueryRow("SELECT nickname FROM users "+
		"WHERE nickname = $1 ", nickname).Scan(&result)
	switch err {
	case pgx.ErrNoRows:
		return "", models.ErrNotExists
	case nil:
		return result, nil
	default:
		return "", err
	}
}
