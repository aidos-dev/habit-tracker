package telegram

import (
	"fmt"

	"github.com/aidos-dev/habit-tracker/internal/clients/telegram"
	"github.com/aidos-dev/habit-tracker/internal/events"
	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/aidos-dev/habit-tracker/internal/repository"
	"github.com/aidos-dev/habit-tracker/pkg/errors"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage repository.Repository
}

type Meta struct {
	ChatID   int
	Username string
}

func New(client *telegram.Client, storage repository.Repository) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, errors.Wrap("can't get events", err)
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

func event(upd models.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: fetchType(upd),
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
