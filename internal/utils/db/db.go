package db

import "github.com/jackc/pgx"

var (
	PgErrUniqueViolation = "23505"
	PgConflict = "P0001"
)

func ErrorCode(err error) string {
	pgerr, ok := err.(pgx.PgError)
	if !ok {
		return ""
	}
	return pgerr.Code
}
