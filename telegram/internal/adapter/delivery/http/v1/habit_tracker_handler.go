package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/aidos-dev/habit-tracker/telegram/internal/models"
	"golang.org/x/exp/slog"
)

func (a *AdapterHandler) UpdateHabitTracker(username string, habitId int, habitTracker models.HabitTracker) {
	const (
		op         = "telegram/internal/adapter/delivery/http/v1/habit_tracker_handler.UpdateHabitTracker"
		habitsUrl  = "/api/habits/"
		trackerUrl = "/tracker"
	)

	a.log.Info(fmt.Sprintf("%s: UpdateHabitTracker method called", op))

	// Perform the necessary logic for command1
	a.log.Info(fmt.Sprintf("%s: Executing UpdateHabitTracker with text: %s", op, username))

	// Make an HTTP request to the backend service
	// http://localhost:8000/telegram/api/habits/7/tracker
	requestURL := a.BackendUrl + habitsUrl + strconv.Itoa(habitId) + trackerUrl + userQuery + username

	type Request struct {
		UnitOfMessure string    `json:"unit_of_messure"`
		Goal          string    `json:"goal"`
		Frequency     string    `json:"frequency"`
		StartDate     time.Time `json:"start_date"`
		EndDate       time.Time `json:"end_date"`
	}

	requestData := Request{
		UnitOfMessure: habitTracker.UnitOfMessure,
		Goal:          habitTracker.Goal,
		Frequency:     habitTracker.Frequency,
		StartDate:     habitTracker.StartDate,
		EndDate:       habitTracker.EndDate,
	}

	a.log.Info(
		fmt.Sprintf("%s: requestData struct prepared for Marshaling", op),
		slog.Any("struct content", requestData),
	)

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		// c.String(http.StatusInternalServerError, err.Error())
		a.log.Error(fmt.Sprintf("%s: failed to encode to JSON", op), sl.Err(err))
		return
	}

	a.log.Info(
		fmt.Sprintf("%s: requestData Marshaling", op),
		slog.Any("Marshaled requestBody content", requestBody),
	)

	// Send a PUT request
	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(requestBody))
	if err != nil {
		// c.String(http.StatusInternalServerError, err.Error())
		a.log.Error(fmt.Sprintf("%s: failed to send http.Put request", op), sl.Err(err))
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		a.log.Error(fmt.Sprintf("%s: failed to execute a request", op), sl.Err(err))
		return
	}

	defer resp.Body.Close()

	response, err := a.readResponse(resp)
	if err != nil {
		a.log.Error(fmt.Sprintf("%s: failed to get the response", op), sl.Err(err))
		return
	}

	// Checking if the "status" field exists in the response
	status, ok := response["status"].(string)
	if !ok {
		a.log.Error(fmt.Sprintf("%s: status not found in response", op))
		return
	}

	a.log.Info(
		fmt.Sprintf("%s: habit tracker has been updated:", op),
		slog.String("status", status),
	)
}

// func (h *Handler) getAllHabitTrackers(c *gin.Context) {
// 	userId, err := getUserId(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	trackers, err := h.services.HabitTracker.GetAll(userId)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getAllHabitTrackers: %v", err.Error()))
// 		return
// 	}

// 	c.JSON(http.StatusOK, trackers)
// }

// func (h *Handler) getHabitTrackerById(c *gin.Context) {
// 	userId, err := getUserId(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	habitId, err := getHabitId(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
// 		return
// 	}

// 	tracker, err := h.services.HabitTracker.GetById(userId, habitId)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getHabitTrackerById: %v", err.Error()))
// 		return
// 	}

// 	c.JSON(http.StatusOK, tracker)
// }

// func (h *Handler) updateHabitTracker(c *gin.Context) {
// 	userId, err := getUserId(c)
// 	if err != nil {
// 		return
// 	}

// 	habitId, err := getHabitId(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
// 		return
// 	}

// 	var input models.UpdateTrackerInput
// 	if err := c.BindJSON(&input); err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	if err := h.services.HabitTracker.Update(userId, habitId, input); err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: updateHabitTracker: %v", err.Error()))
// 		return
// 	}

// 	c.JSON(http.StatusOK, statusResponse{"ok"})
// }

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

/*
method Create temporarily commented out because I decided to create trackers
together with creating a habit. So any tracker is always linked to a certain habit.
This method is not deleted because it might work for some future functionality
*/

// func (h *Handler) createHabitTracker(c *gin.Context) {
// 	_, err := getUserId(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	userHabitId, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
// 		return
// 	}

// 	var input models.HabitTracker
// 	if err := c.BindJSON(&input); err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error from handler: createHabitTracker: %v", err.Error()))
// 		return
// 	}

// 	id, err := h.services.HabitTracker.Create(userHabitId, input)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	c.JSON(http.StatusOK, map[string]interface{}{
// 		"id": id,
// 	})
// }

//////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

/*
method Delete is commented out for the same reasons as method Create
*/

// func (h *Handler) deleteHabitTracker(c *gin.Context) {
// 	userId, err := getUserId(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	habitId, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
// 		return
// 	}

// 	err = h.services.HabitTracker.Delete(userId, habitId)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: deleteHabitTracker: %v", err.Error()))
// 		return
// 	}

// 	c.JSON(http.StatusOK, statusResponse{"ok"})
// }
