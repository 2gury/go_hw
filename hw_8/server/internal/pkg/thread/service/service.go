package service

import (
	"context"
	"log"
	"server/internal/pkg/domain"
	"server/tools"
)

type service struct {
	ThreadRepo domain.ThreadRepository
}

func NewService(threadRepo domain.ThreadRepository) domain.ThreadService {
	return service{
		ThreadRepo: threadRepo,
	}
}

func (s service) Create(ctx context.Context, thread domain.Thread) error {
	requestID, err := tools.GetRequestID(ctx)
	if err != nil {
		return err
	}

	err = s.ThreadRepo.Create(ctx, thread)
	log.Printf("%s %s %v", requestID, "ThreadRepo.Create", err)

	return err
}

func (s service) Get(ctx context.Context, id string) (domain.Thread, error) {
	requestID, err := tools.GetRequestID(ctx)
	if err != nil {
		return domain.Thread{}, err
	}

	thr, err := s.ThreadRepo.Get(ctx, id)
	log.Printf("%s %s %v", requestID, "ThreadRepo.Get", err)

	return thr, err
}
