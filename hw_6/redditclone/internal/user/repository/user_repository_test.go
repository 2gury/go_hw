package repository

import (
	"fmt"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_SelectByUsername(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, userID string, user *models.User)

	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inUsername    string
		outUser       *models.User
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, username string, user *models.User) {
				rows := sqlmock.NewRows([]string{"id", "username", "password"})
				rows.AddRow(user.ID, user.Username, user.Password)
				mock.ExpectQuery(`SELECT`).WithArgs(username).WillReturnRows(rows)
			},
			inUsername: "testuser",
			outUser: &models.User{
				ID:       1,
				Username: "testuser",
				Password: "testpass",
			},
			expError: nil,
		},
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, username string, user *models.User) {
				mock.ExpectQuery(`SELECT`).WithArgs(username).WillReturnError(fmt.Errorf("sql error"))
			},
			inUsername: "testuser",
			expError:   fmt.Errorf("sql error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			userRep := NewUserRepository(db)
			testCase.mockBehaviour(mock, testCase.inUsername, testCase.outUser)
			user, err := userRep.SelectByUsername(testCase.inUsername)

			assert.Equal(t, user, testCase.outUser)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_Insert(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, user *models.User, userID uint64)

	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inUser        *models.User
		outUserID     uint64
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, user *models.User, userID uint64) {
				mock.ExpectBegin()
				lastID := sqlmock.NewRows([]string{"id"}).AddRow(userID)
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs(user.Username, user.Password).
					WillReturnRows(lastID)
				mock.ExpectCommit()
			},
			inUser: &models.User{
				Username: "testuser",
				Password: "testpass",
			},
			outUserID: 1,
			expError:  nil,
		},
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, user *models.User, userID uint64) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs(user.Username, user.Password).
					WillReturnError(fmt.Errorf("sql error"))
				mock.ExpectRollback()
			},
			inUser: &models.User{
				Username: "testuser",
				Password: "testpass",
			},
			outUserID: 0,
			expError:  fmt.Errorf("sql error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			userRep := NewUserRepository(db)
			testCase.mockBehaviour(mock, testCase.inUser, testCase.outUserID)
			userID, err := userRep.Insert(testCase.inUser)

			assert.Equal(t, userID, testCase.outUserID)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}
