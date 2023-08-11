package models

type Channels struct {
	EventCh              chan Event
	StartSendHelloCh     chan bool
	StartSendHelpCh      chan bool
	StartCreateHabitCh   chan bool
	ContinueHabitCh      chan bool
	HabitDataCh          chan Habit
	StartAllHabitsCh     chan bool
	CallAllHabitCh       chan bool
	StartUpdateTrackerCh chan bool
	RequestHabitIdCh     chan bool
	ContinueTrackerCh    chan bool
	ErrChan              chan error
}
