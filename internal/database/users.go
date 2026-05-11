package database

import (
	"context"
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`
}

func (m *UserModel) Create(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	query := "INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id"

	return m.DB.QueryRowContext(ctx, query, user.Email, user.Password, user.Name).Scan(&user.Id)

}

func (m *UserModel) GetUser(query string, args ...interface{}) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := &User{}
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil

}

func (m *UserModel) GetUserByEmail(email string) (*User, error) {
	return m.GetUser("SELECT id, name, email, password FROM users WHERE email = $1", email)
}

func (m *UserModel) GetUserByID(id int) (*User, error) {
	return m.GetUser("SELECT id, name, email, password FROM users WHERE id = $1", id)
}
