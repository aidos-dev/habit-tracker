package telegram

import (
	"errors"
	"fmt"

	"github.com/aidos-dev/habit-tracker/telegram/clients/tgClient"

	"github.com/aidos-dev/habit-tracker/pkg/errs"
	"github.com/aidos-dev/habit-tracker/telegram/events"
	"github.com/aidos-dev/habit-tracker/telegram/models"
	"github.com/aidos-dev/habit-tracker/telegram/storage"
)

type Processor struct {
	tg      *tgClient.Client
	offset  int
	storage storage.Storage
	// handler v1.Handler
	// storage repository.Repository
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func NewProcessor(client *tgClient.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
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
		p.processMessage(event)
	default:
		return errs.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return errs.Wrap("can't process message", err)
	}
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
