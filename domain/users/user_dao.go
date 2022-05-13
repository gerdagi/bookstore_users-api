package users

import (
	"fmt"
	"strings"

	"github.com/gerdagi/bookstore_users-api/datasources/mssql/users_db"
	"github.com/gerdagi/bookstore_users-api/logger"
	resterrors "github.com/gerdagi/bookstore_utils-go/rest_errors"
)

// dao stands for data access object
// only database access sould be in dao not anather place

const (
	queryInsertUser             = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?,?,?,?,?,?);select ID = convert(bigint, SCOPE_IDENTITY());"
	queryGetUser                = "SELECT id, first_name, last_name, email, date_created, status FROM users with (nolock) where id = ?"
	queryUpdateUser             = "UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?"
	queryDeleteUser             = "DELETE from users WHERE id = ?"
	queryFindUserByStatus       = "SELECT id, first_name, last_name, email, date_created, status from users WHERE status = ?"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status from users where email =? and password=? and status = ?"
)

// method da *point kullanmamızın nedeni asıl user üzinde (method da)
// güncelle yapmamız, *pointer kullnamadığımızda copy üzerinde güncelleme yapıyoruz

func (user *User) FindByEmailAndPassword() *resterrors.RestError {
	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email adn password statement", err)
		return resterrors.NewInternalServerError("database error", err)
	}

	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		if strings.Contains(getErr.Error(), "no rows in result set") {
			return resterrors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to get user by email and password", getErr)
		return resterrors.NewInternalServerError("database error", err)
	}

	return nil
}

func (user *User) FindByStatus(status string) ([]User, *resterrors.RestError) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find user statement", err)
		return nil, resterrors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to run find user statement", err)
		return nil, resterrors.NewInternalServerError("database error", err)
	}
	// close sql connection
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		// if we do not send as a point, we send copy of user and after scan we have empty user data
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when trying to run find user statement", err)
			return nil, resterrors.NewInternalServerError("database error", err)
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, resterrors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}

	return results, nil
}

func (user *User) Update() *resterrors.RestError {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return resterrors.NewInternalServerError("database error", err)
	}

	// this is very important, should close statement to close sql connection
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return resterrors.NewInternalServerError("database error", err)
	}

	return nil
}

func (user *User) Delete() *resterrors.RestError {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return resterrors.NewInternalServerError("database error", err)
	}

	// this is very important, should close statement to close sql connection
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)
	if err != nil {
		logger.Error("error when trying to delete user", err)
		return resterrors.NewInternalServerError("database error", err)
	}

	return nil
}

func (user *User) Get() *resterrors.RestError {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return resterrors.NewInternalServerError("database error", err)
	}

	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		logger.Error("error when trying to get user by id", getErr)
		return resterrors.NewInternalServerError("database error", err)
	}

	return nil
}

func (user *User) Save() *resterrors.RestError {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return resterrors.NewInternalServerError("database error", err)
	}

	// this is very important, should close statement to close sql connection
	defer stmt.Close()

	insertResult := stmt.QueryRow(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)

	var userId int32
	if err := insertResult.Scan(&userId); err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return resterrors.NewInternalServerError("database error", err)
	}

	user.Id = userId
	return nil
}
