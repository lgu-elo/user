package user

import (
	"context"

	"github.com/lgu-elo/user/internal/user/model"
	"github.com/pkg/errors"
)

type Repository interface {
	GetById(id int) (*model.User, error)
	GetAll() ([]*model.User, error)
	Update(user *model.User) error
	Delete(id int) error
	Create(user *model.User) error
}

func (db *storage) GetById(id int) (*model.User, error) {
	db.lock.Lock()
	defer db.lock.Unlock()

	var user model.User

	err := db.QueryRow(
		context.Background(),
		"SELECT id,login,name,role FROM users WHERE id=$1",
		id,
	).Scan(&user.ID, &user.Login, &user.Name, &user.Role)

	db.log.Println(user.Login)
	if err != nil {
		return nil, errors.Wrap(err, "user not found")
	}

	return &user, err
}

func (db *storage) GetAll() ([]*model.User, error) {
	var users []*model.User
	db.lock.Lock()
	defer db.lock.Unlock()

	rows, err := db.Query(context.Background(), "SELECT id,login,name,role FROM users")
	if err != nil {
		return nil, errors.Wrap(err, "can't get all users")
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		rows.Scan(&user.ID, &user.Login, &user.Name, &user.Role)

		users = append(users, &user)
	}

	return users, nil
}

func (db *storage) Update(user *model.User) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	rows, err := db.Query(
		context.Background(),
		"UPDATE users set login=$1, name=$2, role=$3 WHERE id=$4",
		user.Login, user.Name, user.Role, user.ID,
	)
	rows.Close()

	if err != nil {
		return errors.Wrap(err, "can't update user")
	}

	if err := rows.Err(); err != nil {
		return errors.Wrap(err, "can't update user")
	}

	return nil
}

func (db *storage) Create(user *model.User) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	rows, err := db.Query(
		context.Background(),
		"INSERT INTO users(login, password, name, role) VALUES($1,$2,$3,$4)",
		user.Login,
		user.Password,
		user.Name,
		user.Role,
	)
	rows.Close()

	if err != nil {
		return errors.Wrap(err, "can't create user")
	}

	if err := rows.Err(); err != nil {
		return errors.Wrap(err, "can't create user")
	}

	return nil
}

func (db *storage) Delete(id int) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	rows, err := db.Query(
		context.Background(),
		"DELETE FROM users WHERE id=$1",
		id,
	)
	rows.Close()

	if err != nil {
		return errors.Wrap(err, "can't delete user")
	}

	if err := rows.Err(); err != nil {
		return errors.Wrap(err, "can't delete user")
	}

	return nil
}
