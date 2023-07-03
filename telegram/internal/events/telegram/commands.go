package telegram

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aidos-dev/habit-tracker/pkg/errs"
	"github.com/aidos-dev/habit-tracker/telegram/internal/models"
	"github.com/aidos-dev/habit-tracker/telegram/internal/storage"
)

const (
	StartCmd    = "/start"
	HelpCmd     = "/help"
	Habit       = "/new_habit"
	DeleteHabit = "/delete_habit"
	Cancel      = "/cancel"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command [%s] from [%s]", text, username)

	// chatIDchan := make(chan int)
	// usernameChan := make(chan string)
	textChan := make(chan string)

	errCh := make(chan error)
	quitCh := make(chan bool)

	// p.wg.Wait()

	switch text {
	case StartCmd:
		errCh <- p.sendHello(chatID, username)
	case HelpCmd:
		errCh <- p.sendHelp(chatID)
	case Habit:
		for {
			textChan <- text
			errCh <- p.createHabit(chatID, username, textChan, quitCh)
			if <-quitCh {
				break
			}
		}

	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}

	err := <-errCh

	log.Printf("doCmd err content is: %v", err)

	return err
}

func (p *Processor) createHabit(chatID int, username string, textChan chan string, quitCh chan bool) (err error) {
	defer func() { err = errs.WrapIfErr("can't do command: save page", err) }()

	log.Print("createHabit method called")

	// defer p.wg.Done()
	// p.wg.Add(1)

	var habit models.Habit
	var tracker models.HabitTracker

	p.mu.Lock()

	switch {
	case habit.Title == "":
		if err := p.tg.SendMessage(chatID, msgHabitTitle); err != nil {
			return err
		}
		text := <-textChan
		if text == Cancel {
			p.mu.Unlock()
			return nil
		}
		habit.Title = text
	case habit.Description == "":
		if err := p.tg.SendMessage(chatID, msgHabitDescription); err != nil {
			return err
		}
		text := <-textChan
		if text == Cancel {
			p.mu.Unlock()
			return nil
		}
		habit.Description = text
	case tracker.UnitOfMessure == "":
		if err := p.tg.SendMessage(chatID, msgUnitOfMessure); err != nil {
			return err
		}
		text := <-textChan
		if text == Cancel {
			p.mu.Unlock()
			return nil
		}
		tracker.UnitOfMessure = text
	case tracker.Frequency == "":
		if err := p.tg.SendMessage(chatID, msgFrequency); err != nil {
			return err
		}
		text := <-textChan
		if text == Cancel {
			p.mu.Unlock()
			return nil
		}
		tracker.Frequency = text
	case tracker.StartDate.IsZero():
		if err := p.tg.SendMessage(chatID, msgStartDate); err != nil {
			return err
		}
		text := <-textChan
		if text == Cancel {
			p.mu.Unlock()
			return nil
		}

		t, err := time.Parse(timeFormat, text)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		tracker.StartDate = t

	case tracker.EndDate.IsZero():
		if err := p.tg.SendMessage(chatID, msgEndDate); err != nil {
			return err
		}
		text := <-textChan
		if text == Cancel {
			p.mu.Unlock()
			return nil
		}

		t, err := time.Parse(timeFormat, text)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		tracker.EndDate = t
		quitCh <- true
	}

	habitId := p.adapter.CreateHabit(username, habit)

	log.Printf("created habit id is: %v", habitId)

	p.adapter.UpdateHabitTracker(username, habitId, tracker)

	// isExists, err := p.storage.IsExists(page)
	// if err != nil {
	// 	return err
	// }
	// if isExists {
	// 	return p.tg.SendMessage(chatID, msgAlreadyExists)
	// }

	// if err := p.storage.Save(page); err != nil {
	// 	return err
	// }

	if err := p.tg.SendMessage(chatID, msgCreated); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int, username string) error {
	log.Print("sendHello method called")

	p.adapter.SignUp(username)

	log.Printf("user [%v] started bot\n", username)

	return p.tg.SendMessage(chatID, msgHello)
}

//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////

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
		return p.tg.SendMessage(chatID, msgHabitAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgCreated); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = errs.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoHabitCreated)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}
