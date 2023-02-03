package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"url-shortener/internal/api"
	database "url-shortener/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/soqet/configjson"
)

type Config struct {
	ApiUrl string `json:"apiUrl"`
	Port   int    `json:"port"`
}

func main() {
	config := new(Config)
	configjson.ReadConfigFile("./config.json", config)
	db, err := database.NewDb("./urls.db")
	if err != nil {
		panic(err)
	}
	router := chi.NewRouter()
	api.Init(router, db, config.ApiUrl)
	srv := http.Server{Addr: fmt.Sprintf(":%d", config.Port), Handler: router}
	log.Println("server started")
	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 1)
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}
	cancel()
	log.Println("server gracefully shutted down")
}
