package models

import (
	"fmt"
	"github.com/VividCortex/mysqlerr"
	"github.com/badThug/otus-social-network/app/customErrors"
	"github.com/badThug/otus-social-network/app/storage"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"strings"
)

const SexMale = 1
const SexFemale = 2
const SexOther = 3

type User struct {
	User_id    int
	Name       string
	Email      string
	Surname    string
	Age        int
	Interests  string
	City       string
	Sex        int
	Password   string
	Created_at string
	Updated_at string
}

func (u *User) getSexTitle() string {
	var sexTitle string
	switch u.Sex {
	case SexMale:
		sexTitle = "male"
	case SexFemale:
		sexTitle = "female"
	default:
		sexTitle = "other"
	}

	return sexTitle
}

func (u *User) Public() map[string]interface{} {
	return map[string]interface{}{
		"user_id":   u.User_id,
		"name":      u.Name,
		"surname":   u.Surname,
		"age":       u.Age,
		"sex":       u.Sex,
		"sexTitle":  u.getSexTitle(),
		"interests": u.Interests,
		"city":      u.City,
		"email":     u.Email,
	}
}

func UserCreate(conn *storage.DbConnection, user *User) (*User, error) {
	db := conn.GetDb()
	insert, err := db.Prepare("INSERT INTO `user`(name, email, password, surname, age, city, interests, sex) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}

	result, err := insert.Exec(user.Name, user.Email, user.Password, user.Surname, user.Age, user.City, user.Interests, user.Sex)
	if driverErr, ok := err.(*mysql.MySQLError); ok {
		if driverErr.Number == mysqlerr.ER_DUP_ENTRY {
			return nil, errors.Wrap(&customErrors.TypedError{fmt.Sprintf("User with email %s already exist", user.Email)}, driverErr.Message)
		}

		return nil, err
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.User_id = int(userId)

	return user, nil
}

func UserFindById(conn *storage.DbConnection, userId int) (*User, error) {
	db := conn.GetDb()

	query := db.QueryRow("SELECT * FROM user WHERE user_id = ?", userId)

	user := &User{}
	err := userQueryScan(query.Scan, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UserFindAllExceptUserId(conn *storage.DbConnection, userId int) ([]*User, error) {
	db := conn.GetDb()

	query, err := db.Query("SELECT * FROM user WHERE user_id != ?", userId)
	if err != nil {
		return nil, err
	}

	collection := []*User{}
	for query.Next() {
		user := &User{}
		err := userQueryScan(query.Scan, user)
		if err != nil {
			return nil, err
		}
		collection = append(collection, user)
	}

	return collection, nil
}

func UserFindByEmail(conn *storage.DbConnection, email string) (*User, error) {
	db := conn.GetDb()

	query := db.QueryRow("SELECT * FROM user WHERE email = ?", email)

	user := &User{}
	err := userQueryScan(query.Scan, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UserFindByUserIds(conn *storage.DbConnection, userIds []int) ([]*User, error) {
	if 0 == len(userIds) {
		return []*User{}, nil
	}

	args := make([]interface{}, len(userIds))
	for i, id := range userIds {
		args[i] = id
	}

	db := conn.GetDb()

	query, err := db.Query("SELECT * FROM user WHERE user_id in (? "+strings.Repeat(",?", len(args)-1)+")", args...)
	if err != nil {
		return nil, err
	}

	collection := []*User{}
	for query.Next() {
		user := &User{}
		err := userQueryScan(query.Scan, user)
		if err != nil {
			return nil, err
		}
		collection = append(collection, user)
	}

	return collection, nil
}

func userQueryScan(scan func(dest ...interface{}) error, user *User) error {
	err := scan(
		&user.User_id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Surname,
		&user.Created_at,
		&user.Updated_at,
		&user.Age,
		&user.City,
		&user.Interests,
		&user.Sex)

	if err != nil {
		return err
	}

	return nil
}