package telegram

import (
	"log"
	"strings"

	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/aidos-dev/habit-tracker/pkg/errs"
)

const (
	HelpCmd = "/help"
	Habit   = "/habit"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command %s from %s", text, username)

	// if isAddCmd(text) {

	// }

	// start: /start: hi + help
	// help: /help
	// habit: /habit

	switch text {
	case HelpCmd:
	case Habit:
	default:
	}
}

func (p *Processor) createHabit(chatID int, text string, username string) (err error) {
	defer func() { err = errs.WrapIfErr("can't do command: create habit", err) }()

	habit := models.Habit{
		Title: text,
	}

	p.handler.
}

// func isAddCmd(text string) bool {

// }
