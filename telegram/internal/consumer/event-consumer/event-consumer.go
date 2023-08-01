package event_consumer

import (
	"fmt"
	"time"

	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/aidos-dev/habit-tracker/telegram/internal/events"
	"golang.org/x/exp/slog"
)

type Consumer struct {
	log       *slog.Logger
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func NewConsumer(log *slog.Logger, fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		log:       log,
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	const op = "telegram/internal/consumer/event-consumer/event-consumer.Start"

	// log.Print("event consumer started")
	c.log.Info(fmt.Sprintf("%s: event consumer started", op))

	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			// log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			// log.Print(err)
			c.log.Error(fmt.Sprintf("%s: failed to handle an event", op), sl.Err(err))

			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	const op = "telegram/internal/consumer/event-consumer/event-consumer.handleEvents"

	for _, event := range events {
		// log.Printf("got new event: %s", event.Text)
		c.log.Info(
			fmt.Sprintf("%s: got new event", op),
			slog.String("event content", event.Text),
		)

		if err := c.processor.Process(event); err != nil {
			// log.Printf("can't handle event: %s", err.Error())
			c.log.Error(fmt.Sprintf("%s: failed to process an event", op), sl.Err(err))

			continue
		}
	}

	return nil
}
