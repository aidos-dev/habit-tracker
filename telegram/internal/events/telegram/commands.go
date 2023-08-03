package telegram

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aidos-dev/habit-tracker/pkg/errs"
	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/aidos-dev/habit-tracker/telegram/internal/models"
	"github.com/aidos-dev/habit-tracker/telegram/internal/storage"
	"golang.org/x/exp/slog"
)

const (
	StartCmd    = "/start"
	HelpCmd     = "/help"
	Habit       = "/new_habit"
	DeleteHabit = "/delete_habit"
	Cancel      = "/cancel"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	const op = "telegram/internal/events/telegram/commands.doCmd"

	text = strings.TrimSpace(text)

	p.log.Info(
		fmt.Sprintf("%s: got new message", op),
		slog.String("message text", text),
		slog.String("from", username),
		slog.Int("chatId", chatID),
	)

	event := models.Event{
		ChatId:   chatID,
		UserName: username,
		Text:     text,
	}

	switch {
	// case text == Cancel:
	// 	p.log.Info(fmt.Sprintf("%s: switch recieved command to Cancel", op))
	// 	return nil
	case text == StartCmd:
		p.startSendHelloCh <- true
		p.log.Info(fmt.Sprintf("%s: switch sent true to startSendHelloCh", op))
	case text == HelpCmd:
		p.startSendHelpCh <- true
		p.log.Info(fmt.Sprintf("%s: switch sent true to startSendHelpCh", op))
	case text == Habit:
		p.startCreateHabitCh <- true
		p.log.Info(fmt.Sprintf("%s: switch sent true to startCreateHabitCh", op))

	default:
		/*
		   this block of code placed to default and wrapped to "select - case"
		   to make it non blocking.
		*/
		select {
		/*
			this case operates when a user started to create a habit and after each input
			the bot prompts a user for the next input (habit title, description, frequency etc.)
			after each input the CreateHabit method sends signal to p.continueHabitCh
			so the CreateHabit method can be recalled again to continue filling habit fields
			and continue creating a habit
		*/
		case <-p.continueHabitCh:
			p.startCreateHabitCh <- true
			p.log.Info(fmt.Sprintf("%s: switch sent true to startCreateHabitCh", op))
		default:

			p.log.Info(fmt.Sprintf("%s: the message couldn't find it's route", op))
			return nil
		}

	}

	p.log.Info(fmt.Sprintf("%s: switch case made its choise", op))

	p.eventCh <- event

	p.log.Info(fmt.Sprintf("%s: event is sent to eventCh", op))

	err := <-p.errChan

	p.log.Info(fmt.Sprintf("%s: err content is: %v", op, err))

	return err
}

func (p *Processor) CreateHabit() {
	const (
		op       = "telegram/internal/events/telegram/commands.CreateHabit"
		habitErr = "can't do command: save page"
	)

	p.log.Info(fmt.Sprintf("%s: goroutine started", op))

	var habit models.Habit
	var tracker models.HabitTracker

	for {

		<-p.startCreateHabitCh

		p.log.Info(fmt.Sprintf("%s: CreateHabit method called", op))

		p.mu.Lock()
		p.log.Info(fmt.Sprintf("%s: locked event channel", op))

		event := <-p.eventCh

		p.mu.Unlock()
		p.log.Info(fmt.Sprintf("%s: unlocked event channel", op))

		p.log.Info(
			fmt.Sprintf("%s: recieved an event", op),
			slog.Any("event content", event),
		)

		chatID := event.ChatId
		username := event.UserName
		text := event.Text

		if text == Cancel {
			p.log.Info(fmt.Sprintf("%s: recieved command to Cancel", op))
			habit = clearHabit(habit)
			tracker = clearTracker(tracker)

			p.log.Info(
				fmt.Sprintf("%s: habit and tracker intermidate values", op),
				slog.Any("habit value", habit),
				slog.Any("tracker value", tracker),
			)

			p.errChan <- nil
		}

		p.log.Info(
			fmt.Sprintf("%s: habit and tracker intermidate values", op),
			slog.Any("habit value", habit),
			slog.Any("tracker value", tracker),
		)

		switch {
		case text == Habit:
			if err := p.tg.SendMessage(chatID, msgHabitTitle); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
				p.log.Info(
					fmt.Sprintf("%s: err sent to errChan", op),
					slog.Any("err content", errs.Wrap(habitErr, err)),
				)
			}

			p.requestNextPromt()

		case habit.Title == "":
			if err := p.tg.SendMessage(chatID, msgHabitDescription); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
				p.log.Info(
					fmt.Sprintf("%s: err sent to errChan", op),
					slog.Any("err content", errs.Wrap(habitErr, err)),
				)
			}

			habit.Title = text
			p.log.Info(
				fmt.Sprintf("%s: habit Title filled", op),
				slog.String("habit title", habit.Title),
			)

			p.requestNextPromt()

		case habit.Description == "":
			if err := p.tg.SendMessage(chatID, msgUnitOfMessure); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
			}

			habit.Description = text
			p.log.Info(
				fmt.Sprintf("%s: habit Description filled", op),
				slog.String("habit description", habit.Description),
			)

			p.requestNextPromt()

		case tracker.UnitOfMessure == "":
			if err := p.tg.SendMessage(chatID, msgFrequency); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
			}

			tracker.UnitOfMessure = text
			p.log.Info(
				fmt.Sprintf("%s: tracker Unit of Messuer filled", op),
				slog.String("tracker UoM", tracker.UnitOfMessure),
			)

			p.requestNextPromt()

		case tracker.Frequency == "":
			if err := p.tg.SendMessage(chatID, msgStartDate); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
			}
			tracker.Frequency = text
			p.log.Info(
				fmt.Sprintf("%s: tracker Frequency filled", op),
				slog.String("tracker frequency", tracker.Frequency),
			)

			p.requestNextPromt()

		case tracker.StartDate.IsZero():
			if err := p.tg.SendMessage(chatID, msgEndDate); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
			}
			t, err := time.Parse(timeFormat, text)
			if err != nil {
				p.log.Error(fmt.Sprintf("%s: failed to parse time", op), sl.Err(err))
				p.errChan <- errs.Wrap(habitErr, err)
			}
			tracker.StartDate = t
			p.log.Info(
				fmt.Sprintf("%s: tracker start date filled", op),
				slog.Any("tracker start date", tracker.StartDate),
			)

			p.requestNextPromt()

		case tracker.EndDate.IsZero():
			// if err := p.tg.SendMessage(chatID, msgEndDate); err != nil {
			// 	p.errChan <- errs.Wrap(habitErr, err)
			// }

			t, err := time.Parse(timeFormat, text)
			if err != nil {
				p.log.Error(fmt.Sprintf("%s: failed to parse time", op), sl.Err(err))
				p.errChan <- errs.Wrap(habitErr, err)
			}
			tracker.EndDate = t
			p.log.Info(
				fmt.Sprintf("%s: tracker end date filled", op),
				slog.Any("tracker end date", tracker.EndDate),
			)

			p.log.Info(
				fmt.Sprintf("%s: habit and tracker final values", op),
				slog.Any("habit value", habit),
				slog.Any("tracker value", tracker),
			)

			habitId := p.adapter.CreateHabit(username, habit)

			p.adapter.UpdateHabitTracker(username, habitId, tracker)
			// log.Printf("CreateHabit: created habit id is: %v", habitId)

			p.log.Info(
				fmt.Sprintf("%s: habit created", op),
				slog.Int("habitId", habitId),
			)

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

		}

	}
}

// clearHabit resets all habit fields back to zero values
func clearHabit(habit models.Habit) models.Habit {
	habit.Title = ""
	habit.Description = ""
	return habit
}

// clearTracker resets all habit tracker fields back to zero values
func clearTracker(tracker models.HabitTracker) models.HabitTracker {
	tracker.UnitOfMessure = ""
	tracker.Frequency = ""
	tracker.StartDate = time.Time{}
	tracker.EndDate = time.Time{}
	return tracker
}

/*
requestNextPromt sends nil to p.errChan to release doCmd func and let it
accept the next command from a user.
Also it sends a signal to p.continueHabitCh to let
doCmd func know that CreateHabit method is still in process of
creating a habit and it is waiting for the next message from a user
*/
func (p *Processor) requestNextPromt() {
	const op = "telegram/internal/events/telegram/commands.CreateHabit"

	p.errChan <- nil
	p.log.Info(
		fmt.Sprintf("%s: nil sent to errChan", op),
		slog.Any("err content", nil),
	)

	p.continueHabitCh <- true
	p.log.Info(fmt.Sprintf("%s: signal sent to continueHabitCh", op))
}

func (p *Processor) SendHelp() {
	const op = "telegram/internal/events/telegram/commands.SendHelp"

	p.log.Info(fmt.Sprintf("%s: goroutine started", op))

	for {
		<-p.startSendHelpCh
		p.log.Info(fmt.Sprintf("%s: method called", op))

		p.mu.Lock()

		p.log.Info(fmt.Sprintf("%s: locked event channel", op))

		event := <-p.eventCh
		chatID := event.ChatId

		err := p.tg.SendMessage(chatID, msgHelp)

		p.errChan <- err

		p.mu.Unlock()

		p.log.Info(fmt.Sprintf("%s: unlocked event channel", op))

	}
}

func (p *Processor) SendHello() {
	const op = "telegram/internal/events/telegram/commands.SendHello"

	p.log.Info(fmt.Sprintf("%s: goroutine started", op))

	for {

		<-p.startSendHelloCh

		p.log.Info(fmt.Sprintf("%s: method called", op))

		p.mu.Lock()

		p.log.Info(fmt.Sprintf("%s: locked event channel", op))

		event := <-p.eventCh

		p.log.Info(
			fmt.Sprintf("%s: event content", op),
			slog.Any("content", event),
		)

		chatID := event.ChatId
		username := event.UserName

		p.adapter.SignUp(username)

		p.log.Info(
			fmt.Sprintf("%s: user started the bot", op),
			slog.String("username", username),
		)

		err := p.tg.SendMessage(chatID, msgHello)

		p.errChan <- err

		p.mu.Unlock()

		p.log.Info(fmt.Sprintf("%s: unlocked event channel", op))
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
