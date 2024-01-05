package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

const (
	backendURL = "http://habit-tracker:8000/telegram"
	habitsUrl  = "/api/habits"
	trackerUrl = "/tracker"
	userQuery  = "?tgUser="
)

type AdapterHandler struct {
	log    *slog.Logger
	Engine *gin.Engine
	// Router     *gin.RouterGroup

	// EventCh      chan models.Event
	// StartHabitCh chan bool
	// HabitCh      chan models.Habit
	// TrackerCh    chan models.HabitTracker
}

func NewAdapterHandler(log *slog.Logger) *AdapterHandler {
	return &AdapterHandler{
		log:    log,
		Engine: gin.New(),

		// EventCh:      eventCh,
		// StartHabitCh: startHabitCh,
		// HabitCh:      habitCh,
		// TrackerCh:    trackerCh,
	}
}

/*
a.readResponse method takes an http.Response as an argument.
It checks the status code and decodes JSON values to a map of key as a string
and value as interface to get any data type as a value.
*/
func (a *AdapterHandler) readResponse(resp *http.Response) (map[string]interface{}, error) {
	const op = "adapter: readResponse"

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: request failed. status: %d", op, resp.StatusCode)
	}

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to read the responseBody. err: %v", op, err)
	}

	// Parsing the JSON response
	var response map[string]interface{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to decode response body. err: %v", op, err)
	}

	return response, nil
}
