package main

import (
	"context"
	"eCar/config"
	"eCar/internal/handler"
	"eCar/internal/service"
	"eCar/internal/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stM, err := storage.NewDBStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	sr := service.NewService(stM, cfg)
	h := handler.NewHandler(sr, cfg)

	go h.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	log.Println("stopping application")

	cancel()

	err = stM.ShutDown()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("shutting down application")

}
