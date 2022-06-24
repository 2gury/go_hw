package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"lectures-2022-1/08_microservices/99_hw/microservice/service"
)

type LogSystem struct {
	eventSubs map[service.Admin_LoggingServer]interface{}
	statsSubs map[service.Admin_StatisticsServer]service.Stat
	mxStats   *sync.Mutex
	mxEvents  *sync.Mutex
	perms map[string][]string
}

func NewLogSystem(permissions map[string][]string) *LogSystem {
	for idp, val := range permissions {
		for idm, meth := range val {
			permissions[idp][idm] = strings.ReplaceAll(meth, "*", "")
		}
	}
	return &LogSystem{
		eventSubs: map[service.Admin_LoggingServer]interface{}{},
		statsSubs: map[service.Admin_StatisticsServer]service.Stat{},
		mxStats:   &sync.Mutex{},
		mxEvents:  &sync.Mutex{},
		perms: permissions,
	}
}

func (ls *LogSystem) AddEventSubscriber(sub service.Admin_LoggingServer) {
	ls.mxEvents.Lock()
	defer ls.mxEvents.Unlock()
	ls.eventSubs[sub] = ""
}

func (ls *LogSystem) RemoveEventSubscriber(sub service.Admin_LoggingServer) {
	delete(ls.eventSubs, sub)
}

func (ls *LogSystem) SendEvent(event *service.Event) {
	ls.mxEvents.Lock()
	defer ls.mxEvents.Unlock()
	for stream := range ls.eventSubs {
		err := stream.Send(event)
		if err != nil {
			ls.RemoveEventSubscriber(stream)
		}
	}
}

func (ls *LogSystem) AddStatSubscriber(sub service.Admin_StatisticsServer, statInterval int) {
	ls.mxStats.Lock()
	ls.statsSubs[sub] = service.Stat{
		Timestamp:  time.Now().Unix(),
		ByConsumer: map[string]uint64{},
		ByMethod:   map[string]uint64{},
	}
	ls.mxStats.Unlock()

	go func() {
		for {
			ls.mxStats.Lock()
			curStat, ok := ls.statsSubs[sub]
			if !ok {
				ls.mxStats.Unlock()
				return
			}

			err := sub.Send(&curStat)
			if err != nil {
				ls.RemoveStatSubscriber(sub)
				ls.mxStats.Unlock()
				return
			}
			ls.statsSubs[sub] = service.Stat{
				Timestamp:  time.Now().Unix(),
				ByConsumer: map[string]uint64{},
				ByMethod:   map[string]uint64{},
			}
			ls.mxStats.Unlock()

			time.Sleep(time.Duration(statInterval) * time.Second)
		}
	}()
}

func (ls *LogSystem) RemoveStatSubscriber(sub service.Admin_StatisticsServer) {
	delete(ls.statsSubs, sub)
}

func (ls *LogSystem) SendStats(consumer, method string) {
	ls.mxStats.Lock()
	defer ls.mxStats.Unlock()
	for _, stat := range ls.statsSubs {
		stat.Timestamp = time.Now().Unix()
		stat.ByMethod[method] += 1
		stat.ByConsumer[consumer] += 1
	}
}

func (ls *LogSystem) logInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	p, _ := peer.FromContext(ctx)
	md, _ := metadata.FromIncomingContext(ctx)
	consumer := md.Get("consumer")
	method := info.FullMethod
	if len(consumer) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Unknown consumer")
	}

	allowMethods, ok := ls.perms[consumer[0]]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Unknown consumer")
	}

	access := false
	for _, meth := range allowMethods {
		if strings.HasPrefix(method, meth) {
			access = true
			break
		}
	}
	if !access {
		return nil, status.Errorf(codes.Unauthenticated, "Unknown consumer")
	}

	event := &service.Event{
		Timestamp: time.Now().Unix(),
		Consumer:  consumer[0],
		Method:    method,
		Host:      p.Addr.String(),
	}
	ls.SendEvent(event)
	ls.SendStats(consumer[0], method)

	reply, err := handler(ctx, req)

	return reply, err
}

func (ls *LogSystem) streamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	p, _ := peer.FromContext(ss.Context())
	md, _ := metadata.FromIncomingContext(ss.Context())
	consumer := md.Get("consumer")
	method := info.FullMethod
	if len(consumer) == 0 {
		return status.Errorf(codes.Unauthenticated, "Unknown consumer")
	}

	allowMethods, ok := ls.perms[consumer[0]]
	if !ok {
		return status.Errorf(codes.Unauthenticated, "Unknown consumer")
	}

	access := false
	for _, meth := range allowMethods {
		if strings.HasPrefix(method, meth) {
			access = true
			break
		}
	}
	if !access {
		return status.Errorf(codes.Unauthenticated, "Unknown consumer")
	}

	event := &service.Event{
		Timestamp: time.Now().Unix(),
		Consumer:  consumer[0],
		Method:    method,
		Host:      p.Addr.String(),
	}
	ls.SendEvent(event)
	ls.SendStats(consumer[0], method)

	err := handler(srv, ss)

	return err
}

func StartMyMicroservice(ctx context.Context, listenAddres string, ACLData string) error {
	lis, err := net.Listen("tcp", listenAddres[len(listenAddres)-5:])
	if err != nil {
		return err
	}

	var perms map[string][]string
	err = json.Unmarshal([]byte(ACLData), &perms)
	if err != nil {
		defer lis.Close()
		return err
	}

	logSystem := NewLogSystem(perms)

	server := grpc.NewServer(grpc.UnaryInterceptor(logSystem.logInterceptor), grpc.StreamInterceptor(logSystem.streamInterceptor))
	service.RegisterAdminServer(server, NewAdminService(ctx, logSystem))
	service.RegisterBizServer(server, NewBizService(ctx))

	log.Println("starting server at " + listenAddres)
	go func() {
		go server.Serve(lis)

		<-ctx.Done()
		server.Stop()
	}()

	return nil
}
