package users

import (
	"fmt"
	"microservices/bookstore_users-api/datasources/mysql/users_db"
	"microservices/bookstore_users-api/logger"
	"microservices/bookstore_users-api/utils/errors"
)

const (
	queryInsertUser       = "insert into users(first_name, last_name, email, date_created, status, password) values(?, ?, ?, ?, ?, ?);"
	queryGetUser          = "select id, first_name, last_name, email, date_created, status, password from users where id = ?;"
	queryUpdateUser       = "update users set first_name = ?, last_name = ?, email = ? where id = ?;"
	queryDeleteUser       = "delete from users where id = ?;"
	queryFindUserByStatus = "select id, first_name, last_name, email, date_created, status, password from users where status = ?;"
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()
	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status, &user.Password); getErr != nil {
		logger.Error("error when trying to get user by id", getErr)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, saveErr := users_db.Client.Prepare(queryInsertUser)
	if saveErr != nil {
		logger.Error("error when trying to prepare save user statement", saveErr)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()
	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)
		return errors.NewInternalServerError("database error")
	}
	userID, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return errors.NewInternalServerError("database error")
	}
	user.Id = userID
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()
	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when trying to delete user", err)
		return errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) FindByStatus(status string) (Users, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find user by status statement", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer stmt.Close()
	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find user by status", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer rows.Close()
	results := make(Users, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.DateCreated,
			&user.Status,
			&user.Password); err != nil {
			logger.Error("error when trying to scan row into user struct user", err)
			return nil, errors.NewInternalServerError("database error")
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		// return nil, errors.NewInternalServerError("database error")
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}
