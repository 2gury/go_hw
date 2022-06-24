package repository

import (
	"context"
	"database/sql"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/user"
	"log"
)

type UserPgRepository struct {
	dbConn *sql.DB
}

func NewUserRepository(conn *sql.DB) user.UserPgRepository {
	return &UserPgRepository{
		dbConn: conn,
	}
}

func (r *UserPgRepository) SelectByUsername(username string) (*models.User, error) {
	usr := &models.User{}

	err := r.dbConn.QueryRow(
		`SELECT id, username, password FROM users
                WHERE username = $1`, username).
		Scan(&usr.ID, &usr.Username, &usr.Password)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (r *UserPgRepository) Insert(usr *models.User) (uint64, error) {
	tx, err := r.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return 0, err
	}
	var lastID int64
	err = tx.QueryRow(
		`INSERT INTO users(username, password)
			    VALUES ($1, $2) RETURNING id`,
		usr.Username, usr.Password).Scan(&lastID)
	if err != nil {
		if rollBackError := tx.Rollback(); rollBackError != nil {
			log.Fatal(rollBackError.Error())
		}
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return uint64(lastID), nil
}
