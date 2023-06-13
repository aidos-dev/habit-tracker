package main

import (
	"github.com/aidos-dev/habit-tracker/telegram/clients/tgClient"
	"github.com/aidos-dev/habit-tracker/telegram/events/telegram"
)

func main() {
	tgClient := tgClient.NewClient(models.TgBotHost, MustToken())

	tgProcessor := telegram.NewProcessor(&tgClient, handlers)
}
