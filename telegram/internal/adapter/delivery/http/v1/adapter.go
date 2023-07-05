package v1

import (
	"github.com/gin-gonic/gin"
)

const backendURL = "http://habit-tracker:8000/telegram"

type AdapterHandler struct {
	Engine *gin.Engine
	// Router     *gin.RouterGroup
	BackendUrl string
	// EventCh      chan models.Event
	// StartHabitCh chan bool
	// HabitCh      chan models.Habit
	// TrackerCh    chan models.HabitTracker
}

func NewAdapterHandler() *AdapterHandler {
	return &AdapterHandler{
		Engine:     gin.New(),
		BackendUrl: backendURL,
		// EventCh:      eventCh,
		// StartHabitCh: startHabitCh,
		// HabitCh:      habitCh,
		// TrackerCh:    trackerCh,
	}
}
