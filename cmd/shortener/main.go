package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/olenka-91/shorturl/internal/handlers"
	"github.com/olenka-91/shorturl/internal/models"
	"github.com/olenka-91/shorturl/internal/service"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Info("Creating services...")
	serv := service.NewService()
	log.Debug("Services created successfully")

	log.Info("Creating handlers...")
	handl := handlers.NewHandler(serv)
	log.Debug("Handlers created successfully")

	log.Info("Creating server...")
	server := new(models.Server)
	log.Debug("Server created successfully")

	go func() {
		log.Info("Starting the HTTP server...")
		//	if err := server.Run(os.Getenv("APP_PORT"), handl.InitRoutes(), db); err != nil {
		if err := server.Run(":8080", handl.InitRoutes()); err != nil {
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
