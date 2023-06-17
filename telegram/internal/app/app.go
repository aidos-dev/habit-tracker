package app

import (
	"log"

	"github.com/aidos-dev/habit-tracker/telegram/config"
	"github.com/aidos-dev/habit-tracker/telegram/internal/clients/tgClient"
	event_consumer "github.com/aidos-dev/habit-tracker/telegram/internal/consumer/event-consumer"
	"github.com/aidos-dev/habit-tracker/telegram/internal/events/telegram"
	"github.com/aidos-dev/habit-tracker/telegram/internal/models"
	"github.com/aidos-dev/habit-tracker/telegram/internal/storage/files"
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

	// fetcher

	// processor
	eventsProcessor := telegram.NewProcessor(tgClient, storage)

	log.Print("service started")

	// consumer.Start(fetcher, processor)

	consumer := event_consumer.NewConsumer(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
