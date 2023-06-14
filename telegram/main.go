package main

import (
	"log"

	"github.com/aidos-dev/habit-tracker/telegram/clients/tgClient"
	"github.com/aidos-dev/habit-tracker/telegram/config"
	event_consumer "github.com/aidos-dev/habit-tracker/telegram/consumer/event-consumer"
	"github.com/aidos-dev/habit-tracker/telegram/events/telegram"
	"github.com/aidos-dev/habit-tracker/telegram/models"
	"github.com/aidos-dev/habit-tracker/telegram/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
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
