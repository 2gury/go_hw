package models

type Post struct {
	ID               string    `json:"id"`
	Author           User      `json:"author"`
	Category         string    `json:"category"`
	Comments         []Comment `json:"comments"`
	CreatedAt        string    `json:"created"`
	Score            int64     `json:"score"`
	Text             string    `json:"text"`
	Title            string    `json:"title"`
	Type             string    `json:"type"`
	UpvotePercentage float64   `json:"upvotePercentage"`
	URL              string    `json:"url"`
	Views            uint64    `json:"views"`
	Votes            []Vote    `json:"votes"`
}
