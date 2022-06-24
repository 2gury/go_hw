//nolint:typecheck
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"lectures-2022-1/08_microservices/99_hw/microservice/service"
)

const (
	// какой адрес-порт слушать серверу
	listenAddr string = "127.0.0.1:8082"

	// кого по каким методам пускать
	ACLData string = `{
	"logger":    ["/service.Admin/Logging"],
	"stat":      ["/service.Admin/Statistics"],
	"biz_user":  ["/service.Biz/Check", "/service.Biz/Add"],
	"biz_admin": ["/service.Biz/*"]
}`
)

// чтобы не было сюрпризов когда где-то не успела преключиться горутина и не успело что-то стортовать
func wait(amout int) {
	time.Sleep(time.Duration(amout) * 10 * time.Millisecond)
}

// утилитарная функция для коннекта к серверу
func getGrpcConn(t *testing.T) *grpc.ClientConn {
	grcpConn, err := grpc.Dial(
		listenAddr,
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("cant connect to grpc: %v", err)
	}
	return grcpConn
}

// получаем контекст с нужнымы метаданными для ACL
func getConsumerCtx(consumerName string) context.Context {
	// ctx, _ := context.WithTimeout(context.Background(), time.Second)
	ctx := context.Background()
	md := metadata.Pairs(
		"consumer", consumerName,
	)
	return metadata.NewOutgoingContext(ctx, md)
}

// старт-стоп сервера
func TestServerStartStop(t *testing.T) {
	ctx, finish := context.WithCancel(context.Background())
	err := StartMyMicroservice(ctx, listenAddr, ACLData)
	if err != nil {
		t.Fatalf("cant start server initial: %v", err)
	}
	wait(1)
	finish() // при вызове этой функции ваш сервер должен остановиться и освободить порт
	wait(1)

	// теперь проверим что вы освободили порт и мы можем стартовать сервер ещё раз
	ctx, finish = context.WithCancel(context.Background())
	err = StartMyMicroservice(ctx, listenAddr, ACLData)
	if err != nil {
		t.Fatalf("cant start server again: %v", err)
	}
	wait(1)
	finish()
	wait(1)
}

// у вас наверняка будет что-то выполняться в отдельных горутинах
// этим тестом мы проверяем что вы останавливаете все горутины которые у вас были и нет утечек
// некоторый запас ( goroutinesPerTwoIterations*5 ) остаётся на случай рантайм горутин
func TestServerLeak(t *testing.T) {
	//return
	goroutinesStart := runtime.NumGoroutine()
	TestServerStartStop(t)
	goroutinesPerTwoIterations := runtime.NumGoroutine() - goroutinesStart

	goroutinesStart = runtime.NumGoroutine()
	goroutinesStat := []int{}
	for i := 0; i <= 25; i++ {
		TestServerStartStop(t)
		goroutinesStat = append(goroutinesStat, runtime.NumGoroutine())
	}
	goroutinesPerFiftyIterations := runtime.NumGoroutine() - goroutinesStart
	if goroutinesPerFiftyIterations > goroutinesPerTwoIterations*5 {
		t.Fatalf("looks like you have goroutines leak: %+v", goroutinesStat)
	}
}

// ACL (права на методы доступа) парсится корректно
func TestACLParseError(t *testing.T) {
	// finish'а тут нет потому что стартовать у вас ничего не должно если не получилось распаковать ACL
	err := StartMyMicroservice(context.Background(), listenAddr, "{.;")
	if err == nil {
		t.Fatalf("expacted error on bad acl json, have nil")
	}
}

//ACL (права на методы доступа) работает корректно
func TestACL(t *testing.T) {
	wait(1)
	ctx, finish := context.WithCancel(context.Background())
	err := StartMyMicroservice(ctx, listenAddr, ACLData)
	if err != nil {
		t.Fatalf("cant start server initial: %v", err)
	}
	wait(1)
	defer func() {
		finish()
		wait(1)
	}()

	conn := getGrpcConn(t)
	defer conn.Close()

	biz := service.NewBizClient(conn)
	adm := service.NewAdminClient(conn)

	for idx, ctx := range []context.Context{
		context.Background(),       // нет поля для ACL
		getConsumerCtx("unknown"),  // поле есть, неизвестный консюмер
		getConsumerCtx("biz_user"), // поле есть, нет доступа
	} {
		_, err = biz.Test(ctx, &service.Nothing{})
		if err == nil {
			t.Fatalf("[%d] ACL fail: expected err on disallowed method", idx)
		} else if code := grpc.Code(err); code != codes.Unauthenticated {
			t.Fatalf("[%d] ACL fail: expected Unauthenticated code, got %v", idx, code)
		}
	}

	// есть доступ
	_, err = biz.Check(getConsumerCtx("biz_user"), &service.Nothing{})
	if err != nil {
		t.Fatalf("ACL fail: unexpected error: %v", err)
	}
	_, err = biz.Check(getConsumerCtx("biz_admin"), &service.Nothing{})
	if err != nil {
		t.Fatalf("ACL fail: unexpected error: %v", err)
	}
	_, err = biz.Test(getConsumerCtx("biz_admin"), &service.Nothing{})
	if err != nil {
		t.Fatalf("ACL fail: unexpected error: %v", err)
	}

	// ACL на методах, которые возвращают поток данных
	logger, err := adm.Logging(getConsumerCtx("unknown"), &service.Nothing{})
	_, err = logger.Recv()
	if err == nil {
		t.Fatalf("ACL fail: expected err on disallowed method")
	} else if code := grpc.Code(err); code != codes.Unauthenticated {
		t.Fatalf("ACL fail: expected Unauthenticated code, got %v", code)
	}
}

func TestLogging(t *testing.T) {
	ctx, finish := context.WithCancel(context.Background())
	err := StartMyMicroservice(ctx, listenAddr, ACLData)
	if err != nil {
		t.Fatalf("cant start server initial: %v", err)
	}
	wait(1)
	defer func() {
		finish()
		wait(1)
	}()

	conn := getGrpcConn(t)
	defer conn.Close()

	biz := service.NewBizClient(conn)
	adm := service.NewAdminClient(conn)

	logStream1, err := adm.Logging(getConsumerCtx("logger"), &service.Nothing{})
	time.Sleep(1 * time.Millisecond)

	logStream2, err := adm.Logging(getConsumerCtx("logger"), &service.Nothing{})

	logData1 := []*service.Event{}
	logData2 := []*service.Event{}

	wait(1)

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-time.After(3 * time.Second):
			fmt.Println("looks like you dont send anything to log stream in 3 sec")
			t.Errorf("looks like you dont send anything to log stream in 3 sec")
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 4; i++ {
			evt, err := logStream1.Recv()
			// log.Println("logger 1", evt, err)
			if err != nil {
				t.Errorf("unexpected error: %v, awaiting event", err)
				return
			}

			if !strings.HasPrefix(evt.GetHost(), "127.0.0.1:") || evt.GetHost() == listenAddr {
				t.Errorf("bad host: %v", evt.GetHost())
				return
			}
			// это грязный хак
			// protobuf добавляет к структуре свои поля, которвые не видны при приведении к строке и при reflect.DeepEqual
			// поэтому берем не оригинал сообщения, а только нужные значения
			logData1 = append(logData1, &service.Event{Consumer: evt.Consumer, Method: evt.Method})
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			evt, err := logStream2.Recv()
			// log.Println("logger 2", evt, err)
			if err != nil {
				t.Errorf("unexpected error: %v, awaiting event", err)
				return
			}
			if !strings.HasPrefix(evt.GetHost(), "127.0.0.1:") || evt.GetHost() == listenAddr {
				t.Errorf("bad host: %v", evt.GetHost())
				return
			}
			// это грязный хак
			// protobuf добавляет к структуре свои поля, которвые не видны при приведении к строке и при reflect.DeepEqual
			// поэтому берем не оригинал сообщения, а только нужные значения
			logData2 = append(logData2, &service.Event{Consumer: evt.Consumer, Method: evt.Method})
		}
	}()

	biz.Check(getConsumerCtx("biz_user"), &service.Nothing{})
	time.Sleep(2 * time.Millisecond)

	biz.Check(getConsumerCtx("biz_admin"), &service.Nothing{})
	time.Sleep(2 * time.Millisecond)

	biz.Test(getConsumerCtx("biz_admin"), &service.Nothing{})
	time.Sleep(2 * time.Millisecond)

	wg.Wait()

	expectedLogData1 := []*service.Event{
		{Consumer: "logger", Method: "/service.Admin/Logging"},
		{Consumer: "biz_user", Method: "/service.Biz/Check"},
		{Consumer: "biz_admin", Method: "/service.Biz/Check"},
		{Consumer: "biz_admin", Method: "/service.Biz/Test"},
	}
	expectedLogData2 := []*service.Event{
		{Consumer: "biz_user", Method: "/service.Biz/Check"},
		{Consumer: "biz_admin", Method: "/service.Biz/Check"},
		{Consumer: "biz_admin", Method: "/service.Biz/Test"},
	}

	if !reflect.DeepEqual(logData1, expectedLogData1) {
		t.Fatalf("logs1 dont match\nhave %+v\nwant %+v", logData1, expectedLogData1)
	}
	if !reflect.DeepEqual(logData2, expectedLogData2) {
		t.Fatalf("logs2 dont match\nhave %+v\nwant %+v", logData2, expectedLogData2)
	}
}

func TestStat(t *testing.T) {
	ctx, finish := context.WithCancel(context.Background())
	err := StartMyMicroservice(ctx, listenAddr, ACLData)
	if err != nil {
		t.Fatalf("cant start server initial: %v", err)
	}
	wait(1)
	defer func() {
		finish()
		wait(2)
	}()

	conn := getGrpcConn(t)
	defer conn.Close()

	biz := service.NewBizClient(conn)
	adm := service.NewAdminClient(conn)

	statStream1, err := adm.Statistics(getConsumerCtx("stat"), &service.StatInterval{IntervalSeconds: 2})
	wait(1)
	statStream2, err := adm.Statistics(getConsumerCtx("stat"), &service.StatInterval{IntervalSeconds: 3})

	mu := &sync.Mutex{}
	stat1 := &service.Stat{}
	stat2 := &service.Stat{}

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for {
			stat, err := statStream1.Recv()
			if err != nil && err != io.EOF {
				// fmt.Printf("unexpected error %v\n", err)
				return
			} else if err == io.EOF {
				break
			}
			// log.Println("stat1", stat, err)
			mu.Lock()
			// это грязный хак
			// protobuf добавляет к структуре свои поля, которвые не видны при приведении к строке и при reflect.DeepEqual
			// поэтому берем не оригинал сообщения, а только нужные значения
			stat1 = &service.Stat{
				ByMethod:   stat.ByMethod,
				ByConsumer: stat.ByConsumer,
			}
			mu.Unlock()
		}
	}()
	go func() {
		for {
			stat, err := statStream2.Recv()
			if err != nil && err != io.EOF {
				// fmt.Printf("unexpected error %v\n", err)
				return
			} else if err == io.EOF {
				break
			}
			// log.Println("stat2", stat, err)
			mu.Lock()
			// это грязный хак
			// protobuf добавляет к структуре свои поля, которвые не видны при приведении к строке и при reflect.DeepEqual
			// поэтому берем не оригинал сообщения, а только нужные значения
			stat2 = &service.Stat{
				ByMethod:   stat.ByMethod,
				ByConsumer: stat.ByConsumer,
			}
			mu.Unlock()
		}
	}()

	wait(1)

	biz.Check(getConsumerCtx("biz_user"), &service.Nothing{})
	biz.Add(getConsumerCtx("biz_user"), &service.Nothing{})
	biz.Test(getConsumerCtx("biz_admin"), &service.Nothing{})

	wait(200) // 2 sec

	expectedStat1 := &service.Stat{
		ByMethod: map[string]uint64{
			"/service.Biz/Check":        1,
			"/service.Biz/Add":          1,
			"/service.Biz/Test":         1,
			"/service.Admin/Statistics": 1,
		},
		ByConsumer: map[string]uint64{
			"biz_user":  2,
			"biz_admin": 1,
			"stat":      1,
		},
	}

	mu.Lock()
	if !reflect.DeepEqual(stat1, expectedStat1) {
		t.Fatalf("stat1-1 dont match\nhave %+v\nwant %+v", stat1, expectedStat1)
	}
	mu.Unlock()

	biz.Add(getConsumerCtx("biz_admin"), &service.Nothing{})

	wait(220) // 2+ sec

	expectedStat1 = &service.Stat{
		Timestamp: 0,
		ByMethod: map[string]uint64{
			"/service.Biz/Add": 1,
		},
		ByConsumer: map[string]uint64{
			"biz_admin": 1,
		},
	}
	expectedStat2 := &service.Stat{
		Timestamp: 0,
		ByMethod: map[string]uint64{
			"/service.Biz/Check": 1,
			"/service.Biz/Add":   2,
			"/service.Biz/Test":  1,
		},
		ByConsumer: map[string]uint64{
			"biz_user":  2,
			"biz_admin": 2,
		},
	}

	mu.Lock()
	if !reflect.DeepEqual(stat1, expectedStat1) {
		t.Fatalf("stat1-2 dont match\nhave %+v\nwant %+v", stat1, expectedStat1)
	}
	if !reflect.DeepEqual(stat2, expectedStat2) {
		t.Fatalf("stat2 dont match\nhave %+v\nwant %+v", stat2, expectedStat2)
	}
	mu.Unlock()

	finish()
}

func __dummyLog() {
	fmt.Println(1)
	log.Println(1)
}