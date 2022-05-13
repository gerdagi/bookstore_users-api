package mssqlutils

import (
	"strings"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/gerdagi/bookstore_users-api/utils/errors"
)

const (
	indexUniqueEmail = "unique index"
	ErrorNoRows      = "no rows in result set"
	queryInsertUser  = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?,?,?,?);"
	queryGetUser     = "SELECT * FROM users with (nolock) where id = ?"
)

func ParseError(err error) *errors.RestError {
	sqlErr, ok := err.(*mssql.Error)
	if !ok {
		if strings.Contains(err.Error(), ErrorNoRows) {
			return errors.NewNotFoundError("no record matching given id")
		}
		return errors.NewInternalServerError("error parsing database response")
	}

	switch sqlErr.Number {
	case 1062:
		return errors.NewBadRequestError("invalid data")
	}

	return errors.NewInternalServerError("error processing request ")
}
