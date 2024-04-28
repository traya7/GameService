package main

import (
	"GameService/domain"
	"GameService/domain/repositories"
	"GameService/handler"
	"GameService/service"
	"log"
	"net/http"
)

type Configs struct {
	BackURL string

	DB_Host string
	DB_Usrn string
	DB_Pswd string
	DB_Name string
}

const (
	username = "root"
	password = "root"
	hostname = "127.0.0.1:3306"
	dbname   = "ecommerce"
)

func main() {

	// TEst
	// sss:= service.NewGameManager()
	// sss.GameLoop()
	//
	// return;
	// MAIN
	cfg := Configs{
		BackURL: "",
	}

	db, err := domain.NewSqlDb()
	if err != nil {
		log.Fatal(err)
	}

	params := handler.Params{
		Psvc: service.NewPlayerService(repositories.NewPlayerRepository(db)),
		Gsvc: service.NewGameService(),
		Csvc: service.NewCallbackService(cfg.BackURL),
	}

	server := http.Server{
		Addr:    ":8001",
		Handler: handler.NewRouter(params),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
