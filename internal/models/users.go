package models

import (
	"database/sql"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Password string
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) GetAll() ([]User, error) {
	rows, err := m.DB.Query("SELECT * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var (
			id       int
			name     string
			username string
			email    string
			password string
		)
		if err := rows.Scan(&id, &name, &username, &email, &password); err != nil {
			return nil, err
		}
		users = append(users, User{ID: id, Name: name, Email: email, Password: password})
	}
	return users, nil
}
