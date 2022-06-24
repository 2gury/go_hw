package service

import (
	"context"
	"log"
	"server/internal/pkg/domain"
	"server/tools"
)

type service struct {
	CommentRepo domain.CommentRepository
	ThreadRepo  domain.ThreadRepository
}

func NewService(commentRepo domain.CommentRepository, threadRepo domain.ThreadRepository) domain.CommentService {
	return service{
		CommentRepo: commentRepo,
		ThreadRepo:  threadRepo,
	}
}

func (s service) Create(ctx context.Context, threadID string, comment domain.Comment) error {
	requestID, err := tools.GetRequestID(ctx)
	if err != nil {
		return err
	}

	if err := s.checkThread(ctx, threadID); err != nil {
		return err
	}
	log.Printf("%s %s %v", requestID, "CommentSvc.checkThread", err)

	err = s.CommentRepo.Create(ctx, comment)
	log.Printf("%s %s %v", requestID, "CommentRepo.Create", err)

	return err
}

func (s service) Like(ctx context.Context, threadID string, commentID string) error {
	if err := s.checkThread(ctx, threadID); err != nil {
		return err
	}

	return s.CommentRepo.Like(ctx, commentID)
}

func (s service) checkThread(ctx context.Context, threadID string) error {
	_, err := s.ThreadRepo.Get(ctx, threadID)

	return err
}
