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
		user.About,
		user.Email,
		user.Fullname,
		user.Nickname,
	)
	return err
}

func (r *UserRepo) Update(user *models.User) error {
	err := r.db.QueryRow("UPDATE users "+
		"SET about = $1, "+
		"email = $2, "+
		"fullname = $3 "+
		"WHERE nickname = $4 " +
		"RETURNING * ",
		&user.About, &user.Email, &user.Fullname, &user.Nickname).Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	
	return err
}

func (r *UserRepo) SelectByNickname(nickname string) (models.User, error) {
	user := models.User{}
	err := r.db.QueryRow("SELECT * FROM users "+
		"WHERE nickname = $1 ", nickname).Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepo) SelectByEmail(email string) (models.User, error) {
	user := models.User{}
	err := r.db.QueryRow("SELECT * FROM users "+
		"WHERE email = $1 ", email).Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
