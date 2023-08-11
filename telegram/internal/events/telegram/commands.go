package telegram

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aidos-dev/habit-tracker/pkg/errs"
	"github.com/aidos-dev/habit-tracker/telegram/internal/models"
	"github.com/aidos-dev/habit-tracker/telegram/internal/storage"
	"golang.org/x/exp/slog"
)

const (
	StartCmd      = "/start"
	HelpCmd       = "/help"
	Habit         = "/new_habit"
	AllHabits     = "/all_habits"
	UpdateTracker = "/update_tracker"
	DeleteHabit   = "/delete_habit"
	Cancel        = "/cancel"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	const op = "telegram/internal/events/telegram/commands.doCmd"

	text = strings.TrimSpace(text)

	p.log.Debug(
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

	case text == StartCmd:
		p.startSendHelloCh <- true
		p.log.Debug(fmt.Sprintf("%s: switch sent true to startSendHelloCh", op))
	case text == HelpCmd:
		p.startSendHelpCh <- true
		p.log.Debug(fmt.Sprintf("%s: switch sent true to startSendHelpCh", op))
	case text == Habit:
		p.startCreateHabitCh <- true
		p.log.Debug(fmt.Sprintf("%s: switch sent true to startCreateHabitCh", op))
	case text == AllHabits:
		p.startAllHabitsCh <- true
		p.log.Debug(fmt.Sprintf("%s: switch sent true to startAllHabitsCh", op))
	case text == UpdateTracker:

		p.startChooseTrackerCh <- true
		p.log.Debug(fmt.Sprintf("%s: switch sent true to startChooseTrackerCh", op))

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
			p.log.Debug(fmt.Sprintf("%s: switch sent true to startCreateHabitCh", op))
		case <-p.receiveHabitIdCh:

			p.log.Debug(fmt.Sprintf("%s: receiveHabitIdCh received a signal", op))

			habitId, err := strconv.Atoi(text)
			if err != nil {
				p.tg.SendMessage(chatID, msgWrongIdFormat)
			}
			habit, err := p.getHabitById(habitId, username)
			if err != nil {
				p.tg.SendMessage(chatID, msgWrongHabitId)
			}

			p.startUpdateTrackerCh <- true
			p.log.Debug(fmt.Sprintf("%s: switch sent true to startUpdateTrackerCh", op))

			p.habitDataCh <- habit
			p.log.Debug(
				fmt.Sprintf("%s: habit is sent to habitDataCh", op),
				slog.Any("habit", habit),
			)

			p.askUnitOfMessure(event.ChatId)

		case <-p.continueTrackerCh:
			p.startUpdateTrackerCh <- true
			p.log.Debug(fmt.Sprintf("%s: switch sent true to startUpdateTrackerCh", op))

		default:

			p.log.Debug(fmt.Sprintf("%s: the message couldn't find it's route", op))
			return nil
		}

	}

	p.log.Debug(fmt.Sprintf("%s: switch case made its choise", op))

	p.eventCh <- event

	p.log.Debug(fmt.Sprintf("%s: event is sent to eventCh", op))

	err := <-p.errChan

	p.log.Debug(fmt.Sprintf("%s: err content is: %v", op, err))

	return err
}

// clearHabit resets all habit fields back to zero values
func clearHabit(habit models.Habit) models.Habit {
	habit.Id = 0
	habit.Title = ""
	habit.Description = ""
	habit.Username = ""
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
Also it sends a signal to "continue" channels to let
doCmd func know that some method is still in process of
creating or updating something and it is waiting for the next message from a user
*/
func (p *Processor) requestNextPromt(nextPromtChan chan bool, chanName string) {
	const op = "telegram/internal/events/telegram/commands.requestNextPromt"

	p.errChan <- nil
	p.log.Info(
		fmt.Sprintf("%s: nil sent to errChan", op),
		slog.Any("err content", nil),
	)

	nextPromtChan <- true
	p.log.Info(fmt.Sprintf("%s: signal sent to %s", op, chanName))
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
