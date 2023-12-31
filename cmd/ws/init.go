package main

import (
	"net/http"
	"time"
    "fmt"
	"os"
	"strconv"
	"github.com/markelca/toggles/flags"
	"github.com/markelca/toggles/storage"
    "github.com/markelca/toggles/cmd/ws/internal/websocket"
)

var (
    APP_PORT       = os.Getenv("APP_PORT")
    REDIS_HOST     = os.Getenv("REDIS_HOST")
    REDIS_PORT_STR = os.Getenv("REDIS_PORT")
    MONGO_HOST     = os.Getenv("MONGO_HOST")
    MONGO_PORT_STR = os.Getenv("MONGO_PORT")
)

func Init() error {
     redisPort, err  := strconv.Atoi(REDIS_PORT_STR)
     if err != nil {
         return err
     }
     mongoPort, err := strconv.Atoi(MONGO_PORT_STR)
     if err != nil {
         return err
     }

    database,err := flags.NewFlagMongoRepository(MONGO_HOST,mongoPort)
    if err != nil {
        return err
    }

    cacheClient := storage.NewRedisClient(REDIS_HOST, redisPort)
    flagService := flags.NewFlagService(cacheClient,database)

    controller := websocket.WSController{FlagService:flagService,CacheClient:cacheClient}

	hub := websocket.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, controller, w, r)
	})

    host := fmt.Sprintf(":%v", APP_PORT)
	server := &http.Server{
		Addr:              host,
		ReadHeaderTimeout: 3 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
        return err
	}

    return nil
}
