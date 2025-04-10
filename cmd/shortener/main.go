package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/olenka-91/shorturl/config"
	"github.com/olenka-91/shorturl/internal/handlers"
	"github.com/olenka-91/shorturl/internal/models"
	"github.com/olenka-91/shorturl/internal/service"
	"github.com/olenka-91/shorturl/internal/storage"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	flag.Parse()

	stor, err := storage.NewStorage(config.Configs)
	if err != nil {
		log.Fatalf("error occured while create data storage: %s", err.Error())
	}

	serv := service.NewService(config.Configs.BaseAddressURL, stor)
	defer serv.CloseDB()

	handl := handlers.NewHandler(serv)

	server := new(models.Server)

	go func() {
		log.Info("Starting the HTTP server...")
		//	if err := server.Run(os.Getenv("APP_PORT"), handl.InitRoutes(), db); err != nil {
		//if err := server.Run(":8080", handl.InitRoutes()); err != nil
		if err := server.Run(config.Configs.ServiceURL, handl.InitRoutes()); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("error occured while running http server: %s", err.Error())
			}
			log.Info("Server stopped running")
		}
	}()

	log.Info("App Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("Shutting down the server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("error occured on server shutting down: %s", err.Error())
	}
}
