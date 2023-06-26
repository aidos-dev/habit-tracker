package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aidos-dev/habit-tracker/telegram/config"
	v1 "github.com/aidos-dev/habit-tracker/telegram/internal/adapter/delivery/http/v1"
	server "github.com/aidos-dev/habit-tracker/telegram/internal/adapter/server/httpServer"
	"github.com/aidos-dev/habit-tracker/telegram/internal/clients/tgClient"
	event_consumer "github.com/aidos-dev/habit-tracker/telegram/internal/consumer/event-consumer"
	"github.com/aidos-dev/habit-tracker/telegram/internal/events/telegram"
	"github.com/aidos-dev/habit-tracker/telegram/internal/models"
	"github.com/aidos-dev/habit-tracker/telegram/internal/storage/files"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func Run() {
	// get telegram token
	telegramToken := config.MustToken()

	// tgClient = telegram.New(token)
	tgClient := tgClient.NewClient(models.TgBotHost, telegramToken)

	storage := files.New(storagePath)

	adapter := v1.NewAdapterHandler()

	srv := new(server.Server)

	ginEng := adapter.Engine

	adapter.Router = ginEng.Group("telegram/")

	go func() {
		if err := srv.Run(viper.GetString("port"), ginEng); err != nil {
			logrus.Printf("error occured while running backend adapter http server: %s", err.Error())
			return
		}
	}()

	// fetcher

	// processor
	eventsProcessor := telegram.NewProcessor(tgClient, storage, adapter)

	log.Print("service started")

	// consumer.Start(fetcher, processor)

	consumer := event_consumer.NewConsumer(eventsProcessor, eventsProcessor, batchSize)

	go func() {
		if err := consumer.Start(); err != nil {
			logrus.Printf("error occured while running telegram consumer service: %s", err.Error())
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("telegram service Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
}
