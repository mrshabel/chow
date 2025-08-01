package main

import (
	"chow/internal/config"
	"chow/internal/database"
	"chow/internal/handler"
	"chow/internal/repository"
	"chow/internal/router"
	"chow/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// get config
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	// db
	db, err := database.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// repos
	userRepo := repository.NewUserRepository(db.DB)

	// services
	authService := service.NewAuthService(cfg, userRepo)

	// handlers
	authHandler := handler.NewAuthHandler(authService)

	// create server router
	r := gin.Default()
	router.RegisterRoutes(r, authHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf("localhost:%v", cfg.Port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// create a done channel to signal when the shutdown is complete
	done := make(chan struct{}, 1)

	// run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	log.Printf("Server starting on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http server error: %s", err)
	}

	// wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}

func gracefulShutdown(apiServer *http.Server, done chan struct{}) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v\n", err)
	}

	log.Println("Server exiting")

	// notify the main goroutine that the shutdown is complete
	done <- struct{}{}
}
