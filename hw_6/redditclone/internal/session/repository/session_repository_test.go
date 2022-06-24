package repository

import (
	"encoding/json"
	"fmt"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	"testing"

	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
	type mockBehaviour func(conn *redigomock.Conn, sessValue string, sess []byte, timeExpire int)
	t.Parallel()

	rdConn := redigomock.NewConn()
	defer rdConn.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inSession     *models.Session
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(conn *redigomock.Conn, sessValue string, sess []byte, timeExpire int) {
				conn.Command("SET", sessValue, sess, "EX", timeExpire).Expect("OK")
			},
			inSession: models.NewSession("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"),
			expError:  nil,
		},
		{
			name: "Error: redis error",
			mockBehaviour: func(conn *redigomock.Conn, sessValue string, sess []byte, timeExpire int) {
				conn.Command("SET", sessValue, sess, "EX", timeExpire).Expect("NOT OK")
			},
			inSession: models.NewSession("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"),
			expError:  fmt.Errorf("redis: not OK"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			bytesSess, err := json.Marshal(testCase.inSession)
			if err != nil {
				t.Error(err.Error())
				return
			}
			testCase.mockBehaviour(rdConn, testCase.inSession.Value, bytesSess, testCase.inSession.GetTime())

			sessRep := NewSessionRdRepository(rdConn)
			err = sessRep.Create(testCase.inSession)

			assert.Equal(t, err, testCase.expError)
		})
	}
}

func Test_Get(t *testing.T) {
	type mockBehaviour func(conn *redigomock.Conn, sessValue string, sess []byte)
	t.Parallel()

	rdConn := redigomock.NewConn()
	defer rdConn.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inSession     *models.Session
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(conn *redigomock.Conn, sessValue string, sess []byte) {
				conn.Command("GET", sessValue).Expect(sess)
			},
			inSession: models.NewSession("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"),
			expError:  nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			bytesSess, err := json.Marshal(testCase.inSession)
			if err != nil {
				t.Error(err.Error())
				return
			}
			testCase.mockBehaviour(rdConn, testCase.inSession.Value, bytesSess)

			sessRep := NewSessionRdRepository(rdConn)
			sess, err := sessRep.Get(testCase.inSession.Value)

			assert.Equal(t, sess, testCase.inSession)
			assert.Equal(t, err, testCase.expError)
		})
	}
}

func Test_Delete(t *testing.T) {
	type mockBehaviour func(conn *redigomock.Conn, sessValue string)
	t.Parallel()

	rdConn := redigomock.NewConn()
	defer rdConn.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inSessValue   string
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(conn *redigomock.Conn, sessValue string) {
				conn.Command("DEL", sessValue).Expect([]byte("1"))
			},
			inSessValue: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expError:    nil,
		},
		{
			name: "OK",
			mockBehaviour: func(conn *redigomock.Conn, sessValue string) {
				conn.Command("DEL", sessValue).ExpectError(fmt.Errorf("redis: not OK"))
			},
			inSessValue: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expError:    fmt.Errorf("redis: not OK"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour(rdConn, testCase.inSessValue)

			sessRep := NewSessionRdRepository(rdConn)
			err := sessRep.Delete(testCase.inSessValue)

			assert.Equal(t, err, testCase.expError)
		})
	}
}
