package repository

import (
	"encoding/xml"
	"io/fs"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/internal/models"
)

func TestUserUsecase_SelectUsers(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name       string
		inFilename string
		expUsers   *models.SearchUsers
		expError   error
	}{
		{
			name:       "OK",
			inFilename: "test_data/dataset_test_OK.xml",
			expUsers: &models.SearchUsers{
				XMLName: xml.Name{
					Space: "",
					Local: "root",
				},
				Users: []models.SearchUser{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "row",
						},
						ID:            0,
						GUID:          "1a6fa827-62f1-45f6-b579-aaead2b47169",
						IsActive:      false,
						Balance:       "$2,144.93",
						Picture:       "http://placehold.it/32x32",
						Age:           22,
						EyeColor:      "green",
						FirstName:     "Boyd",
						LastName:      "Wolf",
						Gender:        "male",
						Company:       "HOPELI",
						Email:         "boydwolf@hopeli.com",
						Phone:         "+1 (956) 593-2402",
						Address:       "586 Winthrop Street, Edneyville, Mississippi, 9555",
						About:         "about",
						Registered:    "2017-02-05T06:23:27 -03:00",
						FavoriteFruit: "apple",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "row",
						},
						ID:            1,
						GUID:          "46c06b5e-dd08-4e26-bf85-b15d280e5e07",
						IsActive:      false,
						Balance:       "$2,705.71",
						Picture:       "http://placehold.it/32x32",
						Age:           21,
						EyeColor:      "green",
						FirstName:     "Hilda",
						LastName:      "Mayer",
						Gender:        "female",
						Company:       "QUINTITY",
						Email:         "hildamayer@quintity.com",
						Phone:         "+1 (932) 421-2117",
						Address:       "311 Friel Place, Loyalhanna, Kansas, 6845",
						About:         "lols",
						Registered:    "2016-11-20T04:40:07 -03:00",
						FavoriteFruit: "banana",
					},
				},
			},
			expError: nil,
		},
		{
			name:       "OK",
			inFilename: "test_data/dataset_test_non_exist.xml",
			expUsers:   nil,
			expError: &fs.PathError{
				Op:   "open",
				Path: "test_data/dataset_test_non_exist.xml",
				Err:  syscall.Errno(2),
			},
		},
		{
			name:       "OK",
			inFilename: "test_data/dataset_test_FAIL.xml",
			expUsers:   nil,
			expError: &xml.SyntaxError{
				Msg:  "unexpected EOF",
				Line: 13,
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			userRep := NewUserRepository(testCase.inFilename)

			users, err := userRep.SelectUsers()

			assert.Equal(t, users, testCase.expUsers)
			assert.Equal(t, err, testCase.expError)
		})
	}
}
