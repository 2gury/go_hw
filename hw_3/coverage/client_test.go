package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/internal/user/repository"
	"gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/internal/user/usecase"
)

func TestFindUsers(t *testing.T) {
	t.Parallel()

	token := "qwerty"
	filename := "dataset.xml"

	testTable := []struct {
		name              string
		inSrvAccessToken  string
		inCliAccessToken  string
		inFilename        string
		inSearchRequest   SearchRequest
		expSearchResponse *SearchResponse
		expError          error
	}{
		{
			name:             "OK",
			inSrvAccessToken: token,
			inCliAccessToken: token,
			inFilename:       filename,
			inSearchRequest: SearchRequest{
				Limit:      3,
				Query:      "do",
				OrderField: "Age",
				OrderBy:    -1,
			},
			expSearchResponse: &SearchResponse{
				Users: []User{
					{
						ID:     13,
						Name:   "Whitley Davidson",
						Age:    40,
						About:  "Consectetur dolore anim veniam aliqua deserunt officia eu. Et ullamco commodo ad officia duis ex incididunt proident consequat nostrud proident quis tempor. Sunt magna ad excepteur eu sint aliqua eiusmod deserunt proident. Do labore est dolore voluptate ullamco est dolore excepteur magna duis quis. Quis laborum deserunt ipsum velit occaecat est laborum enim aute. Officia dolore sit voluptate quis mollit veniam. Laborum nisi ullamco nisi sit nulla cillum et id nisi.\n",
						Gender: "male",
					},
					{
						ID:     32,
						Name:   "Christy Knapp",
						Age:    40,
						About:  "Incididunt culpa dolore laborum cupidatat consequat. Aliquip cupidatat pariatur sit consectetur laboris labore anim labore. Est sint ut ipsum dolor ipsum nisi tempor in tempor aliqua. Aliquip labore cillum est consequat anim officia non reprehenderit ex duis elit. Amet aliqua eu ad velit incididunt ad ut magna. Culpa dolore qui anim consequat commodo aute.\n",
						Gender: "female",
					},
					{
						ID:     6,
						Name:   "Jennings Mays",
						Age:    39,
						About:  "Veniam consectetur non non aliquip exercitation quis qui. Aliquip duis ut ad commodo consequat ipsum cupidatat id anim voluptate deserunt enim laboris. Sunt nostrud voluptate do est tempor esse anim pariatur. Ea do amet Lorem in mollit ipsum irure Lorem exercitation. Exercitation deserunt adipisicing nulla aute ex amet sint tempor incididunt magna. Quis et consectetur dolor nulla reprehenderit culpa laboris voluptate ut mollit. Qui ipsum nisi ullamco sit exercitation nisi magna fugiat anim consectetur officia.\n",
						Gender: "male",
					},
				},
				NextPage: true,
			},
			expError: nil,
		},
		{
			name:             "OK",
			inSrvAccessToken: token,
			inCliAccessToken: token,
			inFilename:       filename,
			inSearchRequest: SearchRequest{
				Limit:      3,
				Query:      "kek",
				OrderField: "Age",
				OrderBy:    -1,
			},
			expSearchResponse: &SearchResponse{
				Users:    []User{},
				NextPage: false,
			},
			expError: nil,
		},
		{
			name:             "Error: Bad Limit < 0",
			inSrvAccessToken: token,
			inCliAccessToken: token,
			inFilename:       filename,
			inSearchRequest: SearchRequest{
				Limit: -10,
			},
			expSearchResponse: nil,
			expError:          fmt.Errorf("limit must be > 0"),
		},
		{
			name:             "Error: Bad Limit > 25",
			inSrvAccessToken: token,
			inCliAccessToken: token,
			inFilename:       filename,
			inSearchRequest: SearchRequest{
				Limit:  100,
				Offset: -10,
			},
			expError: fmt.Errorf("offset must be > 0"),
		},
		{
			name:             "Error: Bad Request",
			inSrvAccessToken: token,
			inCliAccessToken: token,
			inFilename:       "salt" + filename,
			inSearchRequest: SearchRequest{
				OrderField: "kek",
			},
			expError: fmt.Errorf("cant unpack error json: unexpected end of JSON input"),
		},
		{
			name:             "Error: Bad AccessToken",
			inSrvAccessToken: token,
			inCliAccessToken: "salt" + token,
			inFilename:       filename,
			expError:         fmt.Errorf("bad AccessToken"),
		},
		{
			name:             "Error: Bad Filename",
			inSrvAccessToken: token,
			inCliAccessToken: token,
			inFilename:       "salt" + filename,
			expError:         fmt.Errorf("SearchServer fatal error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			rep := repository.NewUserRepository(testCase.inFilename)
			use := usecase.NewUserUsecase(rep)
			hnd := NewUserHandler(use, testCase.inSrvAccessToken)

			testServer := httptest.NewServer(http.HandlerFunc(hnd.SearchServer))
			defer testServer.Close()

			client := SearchClient{
				AccessToken: testCase.inCliAccessToken,
				URL:         testServer.URL,
			}

			response, err := client.FindUsers(testCase.inSearchRequest)
			assert.Equal(t, response, testCase.expSearchResponse)
			assert.Equal(t, err, testCase.expError)

		})
	}
}

func TestFindUsers_MockSearhServer(t *testing.T) {
	type mockBehaviour func(w http.ResponseWriter, r *http.Request)
	t.Parallel()

	testTable := []struct {
		name              string
		mockBehaviour     mockBehaviour
		inURL             string
		inSearchRequest   SearchRequest
		expSearchResponse *SearchResponse
		expError          error
	}{
		{
			name: "Error: OrderField invalid",
			mockBehaviour: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(SearchErrorResponse{Error: `OrderField invalid`}) //nolint:errcheck
			},
			inSearchRequest: SearchRequest{
				OrderField: "kek",
			},
			expSearchResponse: nil,
			expError:          fmt.Errorf("OrderFeld %s invalid", "kek"),
		},
		{
			name: "Error: unknown bad request error",
			mockBehaviour: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(SearchErrorResponse{Error: `error`}) //nolint:errcheck
			},
			expSearchResponse: nil,
			expError:          fmt.Errorf("unknown bad request error: %s", "error"),
		},
		{
			name: "Error: Timeout",
			mockBehaviour: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Second)
			},
			expSearchResponse: nil,
			expError:          fmt.Errorf("timeout for limit=1&offset=0&order_by=0&order_field=&query="),
		},
		{
			name: "Error: unknown error",
			mockBehaviour: func(w http.ResponseWriter, r *http.Request) {
			},
			inURL:             "kek",
			expSearchResponse: nil,
			expError:          fmt.Errorf("unknown error Get \"kek?limit=1&offset=0&order_by=0&order_field=&query=\": unsupported protocol scheme \"\""),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(testCase.mockBehaviour))
			defer testServer.Close()

			client := SearchClient{}
			if testCase.inURL == "" {
				client.URL = testServer.URL
			} else {
				client.URL = testCase.inURL
			}

			response, err := client.FindUsers(testCase.inSearchRequest)
			assert.Equal(t, response, testCase.expSearchResponse)
			assert.Equal(t, err, testCase.expError)

		})
	}
}
