package main

import (
	"github.com/aidos-dev/habit-tracker/telegram/clients/tgClient"
	"github.com/aidos-dev/habit-tracker/telegram/config"
	"github.com/aidos-dev/habit-tracker/telegram/events/telegram"
	"github.com/aidos-dev/habit-tracker/telegram/models"
)

func main() {
	// get telegram token
	telegramToken := config.MustToken()

	// tgClient = telegram.New(token)
	tgClient := tgClient.NewClient(models.TgBotHost, telegramToken)

	// fetcher

	// processor

	// consumer.Start(fetcher, processor)

	tgProcessor := telegram.NewProcessor(&tgClient, handlers)
}
