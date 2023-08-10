package telegram

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aidos-dev/habit-tracker/pkg/errs"
	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
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
	case text == AllHabits:
		p.allHabits(chatID, username)
	case text == UpdateTracker:

		p.chooseTrackerToUpdate(chatID, username)

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
		case <-p.requestHabitIdCh:

			p.log.Info(fmt.Sprintf("%s: ==== requestHabitIdCh received a signal <- =====", op))

			habitId, err := strconv.Atoi(text)
			if err != nil {
				p.tg.SendMessage(chatID, msgWrongIdFormat)
			}
			habit, err := p.getHabitById(habitId, username)
			if err != nil {
				p.tg.SendMessage(chatID, msgWrongHabitId)
			}
			p.habitDataChan <- habit
			p.startUpdateTrackerCh <- true
			p.log.Info(fmt.Sprintf("%s: switch sent true to startUpdateTrackerCh", op))
		case <-p.continueTrackerCh:
			p.startUpdateTrackerCh <- true
			p.log.Info(fmt.Sprintf("%s: switch sent true to startUpdateTrackerCh", op))

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
		habitErr = "can't do command: create habit"
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

		// if text == Cancel {
		// 	p.log.Info(fmt.Sprintf("%s: recieved command to Cancel", op))
		// 	habit = clearHabit(habit)
		// 	tracker = clearTracker(tracker)

		// 	p.log.Info(
		// 		fmt.Sprintf("%s: habit intermidate values", op),
		// 		slog.Any("habit value", habit),
		// 	)

		// 	p.errChan <- nil
		// }

		p.log.Info(
			fmt.Sprintf("%s: habit intermidate values", op),
			slog.Any("habit value", habit),
		)

		switch {

		case text == Cancel:

			p.log.Info(fmt.Sprintf("%s: recieved command to Cancel", op))
			habit = clearHabit(habit)
			tracker = clearTracker(tracker)

			p.log.Info(
				fmt.Sprintf("%s: habit values after /cancel command", op),
				slog.Any("habit value", habit),
			)

			p.errChan <- nil

		case text == Habit:
			if err := p.tg.SendMessage(chatID, msgHabitTitle); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
				p.log.Info(
					fmt.Sprintf("%s: err sent to errChan", op),
					slog.Any("err content", errs.Wrap(habitErr, err)),
				)
			}

			p.requestNextPromt(p.continueHabitCh, "continueHabitCh")

		case habit.Title == "":

			habit.Title = text
			p.log.Info(
				fmt.Sprintf("%s: habit Title filled", op),
				slog.String("habit title", habit.Title),
			)

			if err := p.tg.SendMessage(chatID, msgHabitDescription); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
				p.log.Info(
					fmt.Sprintf("%s: err sent to errChan", op),
					slog.Any("err content", errs.Wrap(habitErr, err)),
				)
			}

			p.requestNextPromt(p.continueHabitCh, "continueHabitCh")

		case habit.Description == "":

			habit.Description = text
			p.log.Info(
				fmt.Sprintf("%s: habit Description filled", op),
				slog.String("habit description", habit.Description),
			)

			p.log.Info(
				fmt.Sprintf("%s: habit final values", op),
				slog.Any("habit value", habit),
			)

			habitId := p.adapter.CreateHabit(username, habit)

			p.log.Info(
				fmt.Sprintf("%s: habit created", op),
				slog.Int("habitId", habitId),
			)

			if err := p.tg.SendMessage(chatID, msgCreated); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
				p.log.Info(
					fmt.Sprintf("%s: err sent to errChan", op),
					slog.Any("err content", errs.Wrap(habitErr, err)),
				)
			}

			if err := p.tg.SendMessage(chatID, msgUnitOfMessure); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
				p.log.Info(
					fmt.Sprintf("%s: err sent to errChan", op),
					slog.Any("err content", errs.Wrap(habitErr, err)),
				)
			}

			/*
			   sending signal to p.requestNextPromt to let know doCmd func
			   that next event should go to UpdateTracker method
			*/
			p.requestNextPromt(p.continueTrackerCh, "continueTrackerCh")

			habitData := models.Habit{
				Id:          habitId,
				Title:       habit.Title,
				Description: habit.Description,
				Username:    username,
			}

			/*
				habitData to be passed to p.habitDataChan
				so UpdateTracker method can use the habit data
				to create or update a tracker for this habit
			*/
			p.habitDataChan <- habitData
			p.log.Info(
				fmt.Sprintf("%s: habitData is sent to p.habitDataChan", op),
			)

			/*
				clean up habit in order to release memory
				and prepare it for other future habits
			*/
			habit = clearHabit(habit)

			p.log.Info(
				fmt.Sprintf("%s: habit values after cleaning up", op),
				slog.Any("habit value", habit),
			)

			/*
				sending signal to p.startUpdateTrackerCh in order to start
				creating a tracker for the habit
			*/
			// p.startUpdateTrackerCh <- true

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

			// err := p.tg.SendMessage(chatID, msgCreated)
			// p.errChan <- err

			// p.errChan <- nil

		}

	}
}

/*
chooseTrackerToUpdate method asks a user to choose a habit
where the tracker needs to be updated. Then sends a signal to
p.requestHabitIdCh
*/
func (p *Processor) chooseTrackerToUpdate(chatId int, username string) {
	const op = "telegram/internal/events/telegram/commands.chooseTrackerToUpdate"

	p.log.Info(fmt.Sprintf("%s: chooseTrackerToUpdate method called", op))

	/*
		here we ask a user to send the ID of a habit where the tracker
		needs to be updated
	*/
	if err := p.tg.SendMessage(chatId, msgChooseHabit); err != nil {
		p.errChan <- err
	}

	p.allHabits(chatId, username)

	/*
		here we send a signal that the next message from a user
		will be a habit Id and it needs to be used
		to find the appropriate habit from the backend
	*/
	p.requestHabitIdCh <- true
}

func (p *Processor) getHabitById(habitId int, username string) (models.Habit, error) {
	const op = "telegram/internal/events/telegram/commands.getHabitById"
	p.log.Info(fmt.Sprintf("%s: getHabitById method called", op))

	var emptyHabit models.Habit
	habit, err := p.adapter.GetHabitById(habitId, username)
	if err != nil {
		return emptyHabit, err
	}

	return habit, err
}

/*
allHabits method gets all habits of a user from the backend db
and sends them to a tg user as a message
*/
func (p *Processor) allHabits(chatId int, username string) {
	const op = "telegram/internal/events/telegram/commands.allHabits"

	p.log.Info(fmt.Sprintf("%s: allHabits method called", op))

	allHabitsData := p.adapter.GetAllHabits(username)

	p.tg.SendMessage(chatId, fmt.Sprintf("%s\n\n%s", msgAllHabits, allHabitsData))
}

func (p *Processor) UpdateTracker() {
	const (
		op       = "telegram/internal/events/telegram/commands.UpdateTracker"
		habitErr = "can't do command: update tracker"
	)

	p.log.Info(fmt.Sprintf("%s: goroutine started", op))

	/*
		habit variable presents here to link created tracker to its parent habit
	*/
	var habit models.Habit

	var tracker models.HabitTracker

	for {

		<-p.startUpdateTrackerCh

		p.log.Info(fmt.Sprintf("%s: UpdateTracker method called", op))

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

		/*
			here the habit is recieved through the p.habitDataChan
			inside the select - case block in order to make this channel non-blocking.
			It is required in this case because here it recieves a habit only in 2 cases:
				1. When a new habit is created it is sent only once to UpdateTracker method (in
				the very beginning of the tracker creation process)
				2. When a user wants to update the tracker for a specific habit, the doCmd will
				send a habit to UpdateTracker method only once (in the very beginning of the tracker
				update process)
		*/
		select {
		case habit = <-p.habitDataChan:
		default:
		}

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

		case text == UpdateTracker:

			if err := p.tg.SendMessage(chatID, msgChooseHabit); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
			}

			p.requestNextPromt(p.continueTrackerCh, "continueTrackerCh")

		case tracker.UnitOfMessure == "":

			tracker.UnitOfMessure = text
			p.log.Info(
				fmt.Sprintf("%s: tracker Unit of Messuer filled", op),
				slog.String("tracker UoM", tracker.UnitOfMessure),
			)

			if err := p.tg.SendMessage(chatID, msgFrequency); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
			}

			p.requestNextPromt(p.continueTrackerCh, "continueTrackerCh")

		case tracker.Frequency == "":

			tracker.Frequency = text
			p.log.Info(
				fmt.Sprintf("%s: tracker Frequency filled", op),
				slog.String("tracker frequency", tracker.Frequency),
			)

			if err := p.tg.SendMessage(chatID, msgStartDate); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
			}

			p.requestNextPromt(p.continueTrackerCh, "continueTrackerCh")

		case tracker.StartDate.IsZero():

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

			if err := p.tg.SendMessage(chatID, msgEndDate); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
			}

			p.requestNextPromt(p.continueTrackerCh, "continueTrackerCh")

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

			p.adapter.UpdateHabitTracker(username, habit.Id, tracker)
			// log.Printf("CreateHabit: created habit id is: %v", habitId)

			p.log.Info(
				fmt.Sprintf("%s: habit created", op),
				slog.Int("habitId", habit.Id),
			)

			/*
				clean up habit and tracker in order to release memory
				and prepare it for other future habits
			*/
			habit = clearHabit(habit)
			tracker = clearTracker(tracker)
			p.log.Info(
				fmt.Sprintf("%s: habit and tracker values after cleaning up", op),
				slog.Any("habit value", habit),
				slog.Any("tracker value", tracker),
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

			err = p.tg.SendMessage(chatID, msgTrackerUpdated)
			p.errChan <- err

		}

	}
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
