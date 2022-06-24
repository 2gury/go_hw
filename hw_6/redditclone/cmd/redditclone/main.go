package main

import (
	"lectures-2022-1/06_databases/99_hw/redditclone/configs/config"
	"lectures-2022-1/06_databases/99_hw/redditclone/configs/server"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/mwares"
	postHnd "lectures-2022-1/06_databases/99_hw/redditclone/internal/post/delivery"
	postRep "lectures-2022-1/06_databases/99_hw/redditclone/internal/post/repository"
	postUse "lectures-2022-1/06_databases/99_hw/redditclone/internal/post/usecases"
	sessRep "lectures-2022-1/06_databases/99_hw/redditclone/internal/session/repository"
	sessUse "lectures-2022-1/06_databases/99_hw/redditclone/internal/session/usecases"
	userHnd "lectures-2022-1/06_databases/99_hw/redditclone/internal/user/delivery"
	userRep "lectures-2022-1/06_databases/99_hw/redditclone/internal/user/repository"
	userUse "lectures-2022-1/06_databases/99_hw/redditclone/internal/user/usecases"
	"log"

	"github.com/gorilla/mux"
)

func main() {
	conf, err := config.LoadConfig("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	postgresConn, err := conf.Postgres.GetPostgresDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer postgresConn.Close()

	redisConn, err := conf.Redis.GetRedisDBConnection()
	if err != nil {
		log.Println(err)
		return
	}
	defer redisConn.Close()

	mx := mux.NewRouter()
	server.ConfigureStatic(mx)

	sessRep := sessRep.NewSessionRdRepository(redisConn)
	sessUse := sessUse.NewSessionUsecase(sessRep)

	mwareManager := mwares.NewMiddlewareManager(sessUse)

	userRep := userRep.NewUserRepository(postgresConn)
	userUse := userUse.NewUserUsecase(userRep)
	userHnd := userHnd.NewUserHandler(userUse, sessUse)
	userHnd.Configure(mx)

	postRep := postRep.NewPostRepository(postgresConn)
	postUse := postUse.NewPostUsecase(postRep)
	postHnd := postHnd.NewPostHandler(postUse)
	postHnd.Configure(mx, mwareManager)

	mx.Use(mwareManager.PanicCoverMiddleware)
	mx.Use(mwareManager.AccessLogMiddleware)

	srv := server.NewHTTPServer(mx, "8080")
	err = srv.Launch()
	if err != nil {
		log.Println(err)
	}
}
