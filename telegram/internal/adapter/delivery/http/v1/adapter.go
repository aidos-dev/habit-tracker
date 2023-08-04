package v1

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

const (
	backendURL = "http://habit-tracker:8000/telegram"
	userQuery  = "?tgUser="
)

type AdapterHandler struct {
	log    *slog.Logger
	Engine *gin.Engine
	// Router     *gin.RouterGroup
	BackendUrl string
	// EventCh      chan models.Event
	// StartHabitCh chan bool
	// HabitCh      chan models.Habit
	// TrackerCh    chan models.HabitTracker
}

func NewAdapterHandler(log *slog.Logger) *AdapterHandler {
	return &AdapterHandler{
		log:        log,
		Engine:     gin.New(),
		BackendUrl: backendURL,
		// EventCh:      eventCh,
		// StartHabitCh: startHabitCh,
		// HabitCh:      habitCh,
		// TrackerCh:    trackerCh,
	}
}
