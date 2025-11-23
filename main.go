package main

import (
	_ "AvitoTestTask/docs/swagger"
	"AvitoTestTask/internal/config"
	"AvitoTestTask/internal/http"
	"AvitoTestTask/internal/repository"
	service "AvitoTestTask/internal/services"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
)

// @title PR Reviewer Assignment Service (Test Task, Fall 2025)
// @version 1.0.0
// @description API for PR Reviewer Assignment Service
// @host localhost:8080
// @BasePath /api/v1

// @tag.name Teams
// @tag.name Users
// @tag.name PullRequests
// @tag.name Health
func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := pgxpool.New(ctx, cfg.BuildDSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepo(db)
	teamRepo := repository.NewTeamRepo(db)
	prRepo := repository.NewPullRequestRepo(db)
	reviewRepo := repository.NewPrReviewsRepo(db)
	statisticRepo := repository.NewStatisticRepo(db)

	userService := service.NewUserServ(userRepo, teamRepo)
	prService := service.NewPullRequestServ(prRepo, reviewRepo, userRepo)
	teamService := service.NewTeamServ(teamRepo, userRepo)
	statsService := service.NewStatisticService(statisticRepo)
	handler := http.NewHandler(&prService, &userService, &teamService, &statsService)

	server := http.NewServer(cfg, handler.Router())
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("server started on port %s", cfg.AppPort)
		serverErrors <- server.Start()
	}()

	shutdown := make(chan os.Signal, 1)

	select {
	case err = <-serverErrors:
		log.Fatalf("server error: %v", err)

	case sig := <-shutdown:
		log.Printf("received signal %v, starting graceful shutdown", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("graceful shutdown failed: %v", err)

			if err := server.Shutdown(ctx); err != nil {
				log.Fatalf("forced shutdown failed: %v", err)
			}
		}

		log.Println("graceful shutdown completed")
	}
}
