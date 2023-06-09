package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	v1 "github.com/aidos-dev/habit-tracker/internal/delivery/http/v1"
	"github.com/aidos-dev/habit-tracker/internal/repository"
	"github.com/aidos-dev/habit-tracker/internal/server"
	"github.com/aidos-dev/habit-tracker/internal/service"
	_ "github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initConfig() error {
	// AddConfigPath receives a derectory name
	viper.AddConfigPath("configs")
	// SetConfig receives a file name (from the directory above)
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func Run() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// get the telegram token
	token := MustToken()

	if err := initConfig(); err != nil {
		logrus.Printf("error occured while running initConfig: %s", err.Error())
		return
	}

	if err := godotenv.Load("build/.env"); err != nil {
		logrus.Printf("error loading env variables: %s", err.Error())
		return
	}

	dbpool, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Printf("failed to initialize db: %s", err.Error())
		return
	}

	repos := repository.NewRepository(dbpool)
	services := service.NewService(repos)
	handlers := v1.NewHandler(services)

	srv := new(server.Server)

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
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

	// if err := dbpool.Close(); err != nil {
	// 	logrus.Errorf("error occured on dbpool connection close: %s", err.Error())
	// }

	defer dbpool.Close()
}
