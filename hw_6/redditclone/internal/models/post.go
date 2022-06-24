package models

import "time"

type Post struct {
	ID               uint64     `json:"id,string"`
	Author           *User      `json:"author"`
	Category         string     `json:"category"`
	CreatedAt        time.Time  `json:"created"`
	Score            int64      `json:"score"`
	Text             string     `json:"text"`
	Title            string     `json:"title"`
	Type             string     `json:"type"`
	UpvotePercentage float64    `json:"upvotePercentage"`
	URL              string     `json:"url"`
	Views            uint64     `json:"views"`
	Comments         []*Comment `json:"comments"`
	Votes            []*Vote    `json:"votes"`

	UserID uint64 `json:"-"`
}
