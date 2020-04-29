package userrepository

import (
	"database/sql"

	"service.auth/models"
)

type UserRepository struct{}

func (u UserRepository) SaveUser(db *sql.DB, user models.User) (models.User, error) {
	stmt := "insert into users (email, password) values($1, $2) RETURNING id;"
	err := db.QueryRow(stmt, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		return user, err
	}
	user.Password = ""
	return user, nil
}

func (u UserRepository) FetchUser(db *sql.DB, user models.User) (models.User, error) {
	row := db.QueryRow("select * from users where email=$1", user.Email)
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}
