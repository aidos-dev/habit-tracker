package telegram

import (
	"errors"
	"fmt"
	"sync"

	v1 "github.com/aidos-dev/habit-tracker/telegram/internal/adapter/delivery/http/v1"
	"github.com/aidos-dev/habit-tracker/telegram/internal/clients/tgClient"
	"golang.org/x/exp/slog"

	"github.com/aidos-dev/habit-tracker/pkg/errs"
	"github.com/aidos-dev/habit-tracker/telegram/internal/events"
	"github.com/aidos-dev/habit-tracker/telegram/internal/models"
	"github.com/aidos-dev/habit-tracker/telegram/internal/storage"
)

type Processor struct {
	log                  *slog.Logger
	tg                   *tgClient.Client
	offset               int
	storage              storage.Storage
	adapter              *v1.AdapterHandler
	mu                   *sync.Mutex
	eventCh              chan models.Event
	startSendHelloCh     chan bool
	startSendHelpCh      chan bool
	startCreateHabitCh   chan bool
	habitDataCh          chan models.Habit
	startAllHabitsCh     chan bool
	startUpdateTrackerCh chan bool
	requestHabitIdCh     chan bool
	continueHabitCh      chan bool
	continueTrackerCh    chan bool
	errChan              chan error
	// HabitCh      chan models.Habit
	// TrackerCh    chan models.HabitTracker
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func NewProcessor(log *slog.Logger, client *tgClient.Client, storage storage.Storage, adapter *v1.AdapterHandler, mu *sync.Mutex, channels models.Channels) *Processor {
	return &Processor{
		log:                  log,
		tg:                   client,
		storage:              storage,
		adapter:              adapter,
		mu:                   mu,
		eventCh:              channels.EventCh,
		startSendHelloCh:     channels.StartSendHelloCh,
		startSendHelpCh:      channels.StartSendHelpCh,
		startCreateHabitCh:   channels.StartCreateHabitCh,
		habitDataCh:          channels.HabitDataCh,
		startAllHabitsCh:     channels.StartAllHabitsCh,
		startUpdateTrackerCh: channels.StartUpdateTrackerCh,
		requestHabitIdCh:     channels.RequestHabitIdCh,
		continueHabitCh:      channels.ContinueHabitCh,
		continueTrackerCh:    channels.ContinueTrackerCh,
		errChan:              channels.ErrChan,
		// HabitCh:      habitCh,
		// TrackerCh:    trackerCh,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	// log.Print("Fetch method called")
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, errs.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no updates found")
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return errs.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	const op = "telegram/internal/events/telegram/telegram.processMessage"

	meta, err := meta(event)
	if err != nil {
		return errs.Wrap("can't process message", err)
	}

	// log.Printf("processMessage: Event content is: [%v]\n", event)
	// the line bellow only for debugging
	p.log.Info(fmt.Sprintf("%s: New event", op), slog.Any("event content", event))

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return errs.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, errs.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd models.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd models.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd models.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
