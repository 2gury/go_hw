package main

import (
	"context"
	"lectures-2022-1/08_microservices/99_hw/microservice/service"
)

type BizService struct {
	ctx context.Context
	service.UnimplementedBizServer
}

func NewBizService(c context.Context) *BizService {
	return &BizService{
		ctx: c,
	}
}

func (bm BizService) Check(context.Context, *service.Nothing) (*service.Nothing, error) {
	return &service.Nothing{}, nil
}

func (bm BizService) Add(context.Context, *service.Nothing) (*service.Nothing, error) {
	return &service.Nothing{}, nil
}

func (bm BizService) Test(context.Context, *service.Nothing) (*service.Nothing, error) {
	return &service.Nothing{}, nil
}
