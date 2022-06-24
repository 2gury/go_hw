package models

type ReceiveMessage struct {
	UpdateID    int         `json:"update_id"`
	ChannelPost ChannelPost `json:"channel_post"`
	Message     Message     `json:"message"`
}

type Message struct {
	MessageID int        `json:"message_id"`
	Date      int        `json:"date"`
	Text      string     `json:"text"`
	Entities  []Entities `json:"entities"`
	From      From       `json:"from"`
	Chat      Chat       `json:"chat"`
}

type ChannelPost struct {
	MessageID int    `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
}

type From struct {
	ID           int    `json:"id"`
	LanguageCode string `json:"language_code"`
	FirstName    string `json:"first_name"`
	UserName     string `json:"username"`
}

type Result struct {
	MessageID int    `json:"message_id"`
	From      From   `json:"from"`
	Chat      Chat   `json:"chat"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
}

type Entities struct {
	Length int    `json:"length"`
	Type   string `json:"type"`
	Offset int    `json:"offset"`
}

type Chat struct {
	ID                          int    `json:"id"`
	Type                        string `json:"type"`
	Title                       string `json:"title"`
	AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
	FirstName                   string `json:"first_name"`
	UserName                    string `json:"username"`
}
