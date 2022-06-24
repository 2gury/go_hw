package domain

import "context"

type Comment struct {
	ID   string
	Text string
}

type CommentService interface {
	Create(context context.Context, threadID string, comment Comment) error
	Like(context context.Context, threadID string, commentID string) error
}

type CommentRepository interface {
	Create(context context.Context, comment Comment) error
	Like(context context.Context, commentID string) error
}
