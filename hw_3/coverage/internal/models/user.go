package models

import (
	"encoding/xml"
	"fmt"
)

type User struct {
	ID     int
	Name   string
	Age    int
	About  string
	Gender string
}

type SearchUser struct {
	XMLName       xml.Name `xml:"row"`
	ID            int      `xml:"id"`
	GUID          string   `xml:"guid"`
	IsActive      bool     `xml:"isActive"`
	Balance       string   `xml:"balance"`
	Picture       string   `xml:"picture"`
	Age           int      `xml:"age"`
	EyeColor      string   `xml:"eyeColor"`
	FirstName     string   `xml:"first_name"`
	LastName      string   `xml:"last_name"`
	Gender        string   `xml:"gender"`
	Company       string   `xml:"company"`
	Email         string   `xml:"email"`
	Phone         string   `xml:"phone"`
	Address       string   `xml:"address"`
	About         string   `xml:"about"`
	Registered    string   `xml:"registered"`
	FavoriteFruit string   `xml:"favoriteFruit"`
}

type SearchUsers struct {
	XMLName xml.Name     `xml:"root"`
	Users   []SearchUser `xml:"row"`
}

func ConverSearchUsersToClientUsers(searchUsers *SearchUsers) []User {
	clienUsers := []User{}
	for _, user := range searchUsers.Users {
		clienUsers = append(clienUsers, User{
			ID:     user.ID,
			Name:   fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			Age:    user.Age,
			About:  user.About,
			Gender: user.Gender,
		})
	}
	return clienUsers
}
