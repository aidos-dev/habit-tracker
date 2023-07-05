package telegram

import (
	"errors"
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

	log.Printf("doCmd: got new command [%s] from [%s]", text, username)

	event := models.Event{
		ChatId:   chatID,
		UserName: username,
		Text:     text,
	}

	switch text {
	case StartCmd:
		p.startSendHelloCh <- true
		log.Print("doCmd: switch sent true to startSendHelloCh")
	case HelpCmd:
		p.startSendHelpCh <- true
	case Habit:
		p.startCreateHabitCh <- true

	}

	log.Print("doCmd: switch case made its choise")

	p.eventCh <- event

	log.Print("doCmd: event is sent to eventCh")

	err := <-p.errChan

	log.Printf("doCmd:  err content is: %v", err)

	return err
}

func (p *Processor) CreateHabit() {
	const habitErr = "can't do command: save page"

	log.Print("CreateHabit: goroutine started")

	for {

		<-p.startCreateHabitCh
		log.Print("CreateHabit: method called")

		p.mu.Lock()
		log.Print("CreateHabit: method locked event channel")

		event := <-p.eventCh

		chatID := event.ChatId
		username := event.UserName
		text := event.Text

		if text == Cancel {
			p.mu.Unlock()
			p.errChan <- nil
		}

		// defer p.wg.Done()
		// p.wg.Add(1)

		var habit models.Habit
		var tracker models.HabitTracker

		switch {
		case habit.Title == "":
			if err := p.tg.SendMessage(chatID, msgHabitTitle); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
			}
			habit.Title = text

		case habit.Description == "":
			if err := p.tg.SendMessage(chatID, msgHabitDescription); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
			}
			habit.Description = text

		case tracker.UnitOfMessure == "":
			if err := p.tg.SendMessage(chatID, msgUnitOfMessure); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
			}
			tracker.UnitOfMessure = text

		case tracker.Frequency == "":
			if err := p.tg.SendMessage(chatID, msgFrequency); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
			}
			tracker.Frequency = text

		case tracker.StartDate.IsZero():
			if err := p.tg.SendMessage(chatID, msgStartDate); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
			}
			t, err := time.Parse(timeFormat, text)
			if err != nil {
				log.Printf("error: CreateHabit: %v", err)
				p.errChan <- errs.Wrap(habitErr, err)
			}
			tracker.StartDate = t

		case tracker.EndDate.IsZero():
			if err := p.tg.SendMessage(chatID, msgEndDate); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
			}

			t, err := time.Parse(timeFormat, text)
			if err != nil {
				log.Printf("error: CreateHabit: %v", err)
				p.errChan <- errs.Wrap(habitErr, err)
			}
			tracker.EndDate = t

			habitId := p.adapter.CreateHabit(username, habit)

			log.Printf("CreateHabit: created habit id is: %v", habitId)

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

			// if err := p.tg.SendMessage(chatID, msgCreated); err != nil {
			// 	p.errChan <- nil
			// }

			err = p.tg.SendMessage(chatID, msgCreated)
			p.errChan <- err

			p.mu.Unlock()
			log.Print("CreateHabit: method unlocked event channel")

		}

	}
}

func (p *Processor) SendHelp() {
	log.Print("SendHelp: goroutine started")

	for {
		<-p.startSendHelpCh

		p.mu.Lock()
		log.Print("SendHelp: method locked event channel")

		event := <-p.eventCh
		chatID := event.ChatId

		err := p.tg.SendMessage(chatID, msgHelp)

		p.errChan <- err

		p.mu.Unlock()
		log.Print("SendHelp: method locked event channel")

	}
}

func (p *Processor) SendHello() {
	log.Print("SendHello: goroutine started")

	for {

		<-p.startSendHelloCh
		log.Print("SendHello: SendHello method called")

		p.mu.Lock()

		log.Print("SendHello: method locked event channel")

		event := <-p.eventCh

		log.Printf("SendHello: event content is: %v", event)

		chatID := event.ChatId
		username := event.UserName

		p.adapter.SignUp(username)

		log.Printf("SendHello: user [%v] started bot\n", username)

		err := p.tg.SendMessage(chatID, msgHello)

		p.errChan <- err

		p.mu.Unlock()

		log.Print("SendHello: method unlocked event channel")
	}
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
