package main

import (
	"lectures-2022-1/05_web_app/99_hw/redditclone/tools/server"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/mwares"
	postHnd "lectures-2022-1/05_web_app/99_hw/redditclone/internal/post/delivery"
	postRep "lectures-2022-1/05_web_app/99_hw/redditclone/internal/post/repository"
	postUse "lectures-2022-1/05_web_app/99_hw/redditclone/internal/post/usecases"
	userHnd "lectures-2022-1/05_web_app/99_hw/redditclone/internal/user/delivery"
	userRep "lectures-2022-1/05_web_app/99_hw/redditclone/internal/user/repository"
	userUse "lectures-2022-1/05_web_app/99_hw/redditclone/internal/user/usecases"
	"log"

	"github.com/gorilla/mux"
)

func main() {
	mx := mux.NewRouter()
	config.ConfigureStatic(mx)

	userRep := userRep.NewUserRepository()
	userUse := userUse.NewUserUsecase(userRep)
	userHnd := userHnd.NewUserHandler(userUse)
	userHnd.Configure(mx)

	postRep := postRep.NewPostRepository()
	postUse := postUse.NewPostUsecase(postRep)
	postHnd := postHnd.NewPostHandler(postUse)
	postHnd.Configure(mx)

	mx.Use(mwares.PanicCoverMiddleware)
	mx.Use(mwares.AccessLogMiddleware)

	srv := config.NewHTTPServer(mx, "8080")
	log.Fatalln(srv.Launch())
}
