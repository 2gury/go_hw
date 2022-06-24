package domain

import "context"

type Thread struct {
	ID   string
	Name string
}

type ThreadService interface {
	Create(ctx context.Context, thread Thread) error
	Get(ctx context.Context, id string) (Thread, error)
}

type ThreadRepository interface {
	Create(ctx context.Context, thread Thread) error
	Get(ctx context.Context, id string) (Thread, error)
}
