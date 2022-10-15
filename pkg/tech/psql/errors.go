package psql

import (
	"errors"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func IsDuplicateKeyErr(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func IsForeignKeyViolationErr(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}

func IsRecordNotFound(err error) bool {
	return errors.Is(gorm.ErrRecordNotFound, err)
}
