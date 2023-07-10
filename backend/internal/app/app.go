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
	"github.com/aidos-dev/habit-tracker/pkg/loggs"
	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	_ "github.com/jackc/pgx/v5"
	"golang.org/x/exp/slog"
)

func Run() {
	// init config: cleanenv
	cfg := config.MustLoad()

	// init logger: slog
	log := loggs.SetupLogger(cfg.Env)

	log.Info(
		"starting habit-tracker with slog logger",
		slog.String("env", cfg.Env),
		slog.String("version", "v1"),
	)
	log.Debug("debug messages are enabled")

	// logrus.SetFormatter(new(logrus.JSONFormatter))

	dbpool, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		// logrus.Printf("failed to initialize db: %s", err.Error())
		log.Error("failed to initialize db", sl.Err(err))
		return
	}

	repos := postgres.NewPostgresRepository(dbpool)
	services := service.NewService(repos)
	handlers := v1.NewHandler(log, services)

	srv := new(server.Server)

	go func() {
		if err := srv.Run(cfg, handlers.InitRoutes()); err != nil {
			// logrus.Printf("error occured while running http server: %s", err.Error())
			log.Error("failed to run http server", sl.Err(err))
			return
		}
	}()

	log.Info("HabbitTrackerApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// logrus.Println("HabbitTrackerApp Shutting Down")
	log.Info("HabbitTrackerApp Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		// logrus.Errorf("error occured on server shutting down: %s", err.Error())
		log.Error("error occured on server shutting down", sl.Err(err))
	}

	log.Info("HabbitTrackerApp Shutting Down 11111222223333333333")

	defer dbpool.Close()
}
