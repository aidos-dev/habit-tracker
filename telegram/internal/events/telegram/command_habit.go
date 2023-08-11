package telegram

import (
	"fmt"

	"github.com/aidos-dev/habit-tracker/pkg/errs"
	"github.com/aidos-dev/habit-tracker/telegram/internal/models"
	"golang.org/x/exp/slog"
)

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

		p.log.Debug(fmt.Sprintf("%s: CreateHabit method called", op))

		p.mu.Lock()
		p.log.Debug(fmt.Sprintf("%s: locked event channel", op))

		event := <-p.eventCh

		p.mu.Unlock()
		p.log.Debug(fmt.Sprintf("%s: unlocked event channel", op))

		p.log.Debug(
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

		p.log.Debug(
			fmt.Sprintf("%s: habit intermidate values", op),
			slog.Any("habit value", habit),
		)

		switch {

		case text == Cancel:

			p.log.Debug(fmt.Sprintf("%s: recieved command to Cancel", op))
			habit = clearHabit(habit)
			tracker = clearTracker(tracker)

			p.log.Debug(
				fmt.Sprintf("%s: habit values after /cancel command", op),
				slog.Any("habit value", habit),
			)

			p.errChan <- nil

		case text == Habit:
			if err := p.tg.SendMessage(chatID, msgHabitTitle); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
				p.log.Debug(
					fmt.Sprintf("%s: err sent to errChan", op),
					slog.Any("err content", errs.Wrap(habitErr, err)),
				)
			}

			p.requestNextPromt(p.continueHabitCh, "continueHabitCh")

		case habit.Title == "":

			habit.Title = text
			p.log.Debug(
				fmt.Sprintf("%s: habit Title filled", op),
				slog.String("habit title", habit.Title),
			)

			if err := p.tg.SendMessage(chatID, msgHabitDescription); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
				p.log.Debug(
					fmt.Sprintf("%s: err sent to errChan", op),
					slog.Any("err content", errs.Wrap(habitErr, err)),
				)
			}

			p.requestNextPromt(p.continueHabitCh, "continueHabitCh")

		case habit.Description == "":

			habit.Description = text
			p.log.Debug(
				fmt.Sprintf("%s: habit Description filled", op),
				slog.String("habit description", habit.Description),
			)

			p.log.Debug(
				fmt.Sprintf("%s: habit final values", op),
				slog.Any("habit value", habit),
			)

			habitId := p.adapter.CreateHabit(username, habit)

			p.log.Debug(
				fmt.Sprintf("%s: habit created", op),
				slog.Int("habitId", habitId),
			)

			if err := p.tg.SendMessage(chatID, msgCreated); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
				p.log.Debug(
					fmt.Sprintf("%s: err sent to errChan", op),
					slog.Any("err content", errs.Wrap(habitErr, err)),
				)
			}

			/*
			   sending signal to p.requestNextPromt to let know doCmd func
			   that next event should go to UpdateTracker method
			*/
			// p.requestNextPromt(p.continueTrackerCh, "continueTrackerCh")

			// p.startUpdateTrackerCh <- true

			// p.eventCh <- event

			p.askUnitOfMessure(event.ChatId)

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
			p.habitDataCh <- habitData
			p.log.Debug(
				fmt.Sprintf("%s: habitData is sent to p.habitDataChan", op),
			)

			/*
				clean up habit in order to release memory
				and prepare it for other future habits
			*/
			habit = clearHabit(habit)

			p.log.Debug(
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

func (p *Processor) getHabitById(habitId int, username string) (models.Habit, error) {
	const op = "telegram/internal/events/telegram/commands.getHabitById"
	p.log.Debug(fmt.Sprintf("%s: getHabitById method called", op))

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
func (p *Processor) AllHabits() {
	const op = "telegram/internal/events/telegram/commands.AllHabits"

	p.log.Info(fmt.Sprintf("%s: goroutine started", op))

	for {

		<-p.startAllHabitsCh

		p.log.Info(fmt.Sprintf("%s: allHabits method called", op))

		event := <-p.eventCh

		allHabitsData := p.adapter.GetAllHabits(event.UserName)

		p.tg.SendMessage(event.ChatId, fmt.Sprintf("%s\n\n%s", msgAllHabits, allHabitsData))

		p.errChan <- nil
	}
}
