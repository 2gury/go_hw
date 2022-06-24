package main

import (
	"context"
	"lectures-2022-1/08_microservices/99_hw/microservice/service"
)

type AdminService struct {
	ctx       context.Context
	logSystem *LogSystem
	service.UnimplementedAdminServer
}

func NewAdminService(c context.Context, ls *LogSystem) *AdminService {
	return &AdminService{
		ctx:       c,
		logSystem: ls,
	}
}

func (am *AdminService) Logging(nothing *service.Nothing, stream service.Admin_LoggingServer) error {
	am.logSystem.AddEventSubscriber(stream)

	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	case <-am.ctx.Done():
		return am.ctx.Err()
	}
}

func (am *AdminService) Statistics(time *service.StatInterval, stream service.Admin_StatisticsServer) error {
	am.logSystem.AddStatSubscriber(stream, int(time.GetIntervalSeconds()))

	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	case <-am.ctx.Done():
		return am.ctx.Err()
	}
}
