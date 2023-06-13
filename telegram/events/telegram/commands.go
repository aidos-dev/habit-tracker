package telegram

import (
	"log"
	"net/url"
	"strings"

	"github.com/aidos-dev/habit-tracker/pkg/errs"
	"github.com/aidos-dev/habit-tracker/telegram/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
	// Habit   = "/habit"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command %s from %s", text, username)

	if isAddCmd(text) {
		// TODO: AddPage()
	}

	// start: /start: hi + help
	// help: /help
	// habit: /habit

	switch text {
	case RndCmd:
	case HelpCmd:
	case StartCmd:
	default:
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = errs.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}

// func (p *Processor) createHabit(chatID int, text string, username string) (err error) {
// 	defer func() { err = errs.WrapIfErr("can't do command: create habit", err) }()

// 	habit := models.Habit{
// 		Title: text,
// 	}

// 	p.handler.
// }

// func isAddCmd(text string) bool {

// }
