package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/aidos-dev/habit-tracker/backend/internal/config"
	v1 "github.com/aidos-dev/habit-tracker/backend/internal/delivery/http/v1"
	"github.com/aidos-dev/habit-tracker/backend/internal/repository/postgres"
	"github.com/aidos-dev/habit-tracker/backend/internal/server"
	"github.com/aidos-dev/habit-tracker/backend/internal/service"
	_ "github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

func Run() {
	// init config: cleanenv
	cfg := config.MustLoad()

	logrus.SetFormatter(new(logrus.JSONFormatter))

	dbpool, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		logrus.Printf("failed to initialize db: %s", err.Error())
		return
	}

	repos := postgres.NewPostgresRepository(dbpool)
	services := service.NewService(repos)
	handlers := v1.NewHandler(services)

	srv := new(server.Server)

	go func() {
		if err := srv.Run(cfg, handlers.InitRoutes()); err != nil {
			logrus.Printf("error occured while running http server: %s", err.Error())
			return
		}
	}()

	logrus.Println("HabbitTrackerApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("HabbitTrackerApp Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	defer dbpool.Close()
}
