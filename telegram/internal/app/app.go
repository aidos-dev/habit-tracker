package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/aidos-dev/habit-tracker/pkg/loggs"
	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/aidos-dev/habit-tracker/telegram/config"
	v1 "github.com/aidos-dev/habit-tracker/telegram/internal/adapter/delivery/http/v1"
	server "github.com/aidos-dev/habit-tracker/telegram/internal/adapter/server/httpServer"
	"github.com/aidos-dev/habit-tracker/telegram/internal/clients/tgClient"
	event_consumer "github.com/aidos-dev/habit-tracker/telegram/internal/consumer/event-consumer"
	"github.com/aidos-dev/habit-tracker/telegram/internal/events/telegram"
	"github.com/aidos-dev/habit-tracker/telegram/internal/models"
	"github.com/aidos-dev/habit-tracker/telegram/internal/storage/files"
	"golang.org/x/exp/slog"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func Run() {
	// init config: cleanenv
	cfg := config.MustLoad()

	// init logger: slog
	log := loggs.SetupLogger(cfg.Env)

	log.Info(
		"starting Telegram service with slog logger",
		slog.String("env", cfg.Env),
		slog.String("version", "v1"),
	)
	log.Debug("debug messages are enabled")

	// get telegram token
	telegramToken := config.MustToken()

	// tgClient = telegram.New(token)
	tgClient := tgClient.NewClient(models.TgBotHost, telegramToken)

	storage := files.New(storagePath)

	var (
		eventCh              = make(chan models.Event)
		startSendHelloCh     = make(chan bool)
		startSendHelpCh      = make(chan bool)
		startCreateHabitCh   = make(chan bool)
		continueHabitCh      = make(chan bool)
		habitDataCh          = make(chan models.Habit)
		startAllHabitsCh     = make(chan bool)
		startUpdateTrackerCh = make(chan bool)
		startChooseTrackerCh = make(chan bool)
		startAskUnitOfMesCh  = make(chan bool)
		receiveHabitIdCh     = make(chan bool)
		continueTrackerCh    = make(chan bool)
		errChan              = make(chan error)
		// habitCh      chan models.Habit
		// trackerCh    chan models.HabitTracker
	)

	channels := models.Channels{
		EventCh:              eventCh,
		StartSendHelloCh:     startSendHelloCh,
		StartSendHelpCh:      startSendHelpCh,
		StartCreateHabitCh:   startCreateHabitCh,
		ContinueHabitCh:      continueHabitCh,
		HabitDataCh:          habitDataCh,
		StartAllHabitsCh:     startAllHabitsCh,
		StartUpdateTrackerCh: startUpdateTrackerCh,
		StartChooseTrackerCh: startChooseTrackerCh,
		StartAskUnitOfMesCh:  startAskUnitOfMesCh,
		ReceiveHabitIdCh:     receiveHabitIdCh,
		ContinueTrackerCh:    continueTrackerCh,
		ErrChan:              errChan,
	}

	adapter := v1.NewAdapterHandler(log)

	srv := new(server.Server)

	ginEng := adapter.Engine

	// adapter.Router = ginEng.Group("/telegram")

	go func() {
		if err := srv.Run(cfg, log, ginEng); err != nil {
			log.Error("error occured while running the adapter to backend http server: %s", sl.Err(err))
			return
		}
	}()

	// fetcher

	// wait group
	// var wg *sync.WaitGroup
	// mutex
	mu := &sync.Mutex{}

	// processor
	eventsProcessor := telegram.NewProcessor(log, tgClient, storage, adapter, mu, channels)

	go eventsProcessor.SendHello()

	go eventsProcessor.SendHelp()

	/*
		method CreateHabit runs in a separate goroutine and keeps listening
		for chanels. This way it can handle a "dialog" with a user while
		a user is in process of habit creation
	*/
	go eventsProcessor.CreateHabit()

	go eventsProcessor.AllHabits()

	/*
		method UpdateTracker runs in a separate goroutine and keeps listening
		for chanels. This way it can handle a "dialog" with a user while
		a user is in process of updating a tracker of the habit
	*/
	go eventsProcessor.UpdateTracker()

	go eventsProcessor.ChooseTrackerToUpdate()

	go eventsProcessor.AskUnitOfMessure()

	// consumer.Start(fetcher, processor)

	consumer := event_consumer.NewConsumer(log, eventsProcessor, eventsProcessor, batchSize)

	go func() {
		if err := consumer.Start(); err != nil {
			log.Error("error occured while running telegram consumer service: %s", sl.Err(err))
			return
		}
	}()

	log.Info("Telegram service has started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("Telegram service Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Error("error occured on server shutting down: %s", sl.Err(err))
	}
}
