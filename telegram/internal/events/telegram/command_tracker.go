package telegram

import (
	"fmt"
	"time"

	"github.com/aidos-dev/habit-tracker/pkg/errs"
	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/aidos-dev/habit-tracker/telegram/internal/models"
	"golang.org/x/exp/slog"
)

/*
chooseTrackerToUpdate method asks a user to choose a habit
where the tracker needs to be updated. Then sends a signal to
p.requestHabitIdCh
*/
func (p *Processor) ChooseTrackerToUpdate() {
	const op = "telegram/internal/events/telegram/commands.ChooseTrackerToUpdate"

	p.log.Info(fmt.Sprintf("%s: goroutine started", op))

	for {
		<-p.startChooseTrackerCh

		p.log.Debug(fmt.Sprintf("%s: ChooseTrackerToUpdate method called", op))

		event := <-p.eventCh

		p.startAllHabitsCh <- true

		p.eventCh <- event

		/*
			here we ask a user to send the ID of a habit where the tracker
			needs to be updated
		*/
		if err := p.tg.SendMessage(event.ChatId, msgChooseHabit); err != nil {
			p.errChan <- err
		}

		p.errChan <- nil

		/*
			here we send a signal that the next message from a user
			will be a habit Id and it needs to be used
			to find the appropriate habit from the backend
		*/
		p.receiveHabitIdCh <- true

	}
}

/*
askUnitOfMessure method logic is in a separate func as it is called in 2
different places:
1) inside the CreateHabit method. It is called when a habit successfully created
2) in a doCmd func when a user wants to update a tracker of an existing habit
*/
func (p *Processor) askUnitOfMessure(chatID int) {
	const op = "telegram/internal/events/telegram/commands.askUnitOfMessure"

	if err := p.tg.SendMessage(chatID, msgUnitOfMessure); err != nil {
		p.errChan <- fmt.Errorf("%s: failed to send a message to a user: %w", op, err)
		p.log.Debug(
			fmt.Sprintf("%s: err sent to errChan", op),
			slog.Any("err content", fmt.Errorf("%s: failed to send a message to a user: %w", op, err)),
		)
	}

	p.requestNextPromt(p.continueTrackerCh, "continueTrackerCh")
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

		p.log.Debug(fmt.Sprintf("%s: UpdateTracker method called", op))

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
		case habit = <-p.habitDataCh:
		default:
		}

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

		if text == Cancel {
			p.log.Debug(fmt.Sprintf("%s: recieved command to Cancel", op))
			habit = clearHabit(habit)
			tracker = clearTracker(tracker)

			p.log.Info(
				fmt.Sprintf("%s: habit and tracker intermidate values", op),
				slog.Any("habit value", habit),
				slog.Any("tracker value", tracker),
			)

			p.errChan <- nil
		}

		p.log.Debug(
			fmt.Sprintf("%s: habit and tracker intermidate values", op),
			slog.Any("habit value", habit),
			slog.Any("tracker value", tracker),
		)

		// if tracker.UnitOfMessure == "" {
		// 	if err := p.tg.SendMessage(chatID, msgUnitOfMessure); err != nil {
		// 		p.errChan <- errs.Wrap(habitErr, err)
		// 		p.log.Debug(
		// 			fmt.Sprintf("%s: err sent to errChan", op),
		// 			slog.Any("err content", errs.Wrap(habitErr, err)),
		// 		)
		// 	}

		// 	p.requestNextPromt(p.continueTrackerCh, "continueTrackerCh")
		// }

		switch {

		case text == UpdateTracker:

			if err := p.tg.SendMessage(chatID, msgChooseHabit); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
			}

			p.requestNextPromt(p.continueTrackerCh, "continueTrackerCh")

		case tracker.UnitOfMessure == "":

			tracker.UnitOfMessure = text
			p.log.Debug(
				fmt.Sprintf("%s: tracker Unit of Messuer filled", op),
				slog.String("tracker UoM", tracker.UnitOfMessure),
			)

			if err := p.tg.SendMessage(chatID, msgFrequency); err != nil {
				p.errChan <- errs.Wrap(habitErr, err)
			}

			p.requestNextPromt(p.continueTrackerCh, "continueTrackerCh")

		case tracker.Frequency == "":

			tracker.Frequency = text
			p.log.Debug(
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
			p.log.Debug(
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
			p.log.Debug(
				fmt.Sprintf("%s: tracker end date filled", op),
				slog.Any("tracker end date", tracker.EndDate),
			)

			p.log.Debug(
				fmt.Sprintf("%s: habit and tracker final values", op),
				slog.Any("habit value", habit),
				slog.Any("tracker value", tracker),
			)

			p.adapter.UpdateHabitTracker(username, habit.Id, tracker)
			// log.Printf("CreateHabit: created habit id is: %v", habitId)

			p.log.Debug(
				fmt.Sprintf("%s: habit created", op),
				slog.Int("habitId", habit.Id),
			)

			/*
				clean up habit and tracker in order to release memory
				and prepare it for other future habits
			*/
			habit = clearHabit(habit)
			tracker = clearTracker(tracker)
			p.log.Debug(
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
