package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/aidos-dev/habit-tracker/telegram/internal/models"
	"golang.org/x/exp/slog"
)

func (a *AdapterHandler) CreateHabit(username string, habit models.Habit) int {
	const (
		op        = "telegram/internal/adapter/delivery/http/v1/habit_handler.CreateHabit"
		habitsUrl = "/api/habits"
	)

	a.log.Info(fmt.Sprintf("%s: CreateHabit method called", op))

	// Perform the necessary logic for command1
	a.log.Info(fmt.Sprintf("%s: Executing CreateHabit with text: %s", op, username))

	// Make an HTTP request to the backend service
	requestURL := a.BackendUrl + habitsUrl + userQuery + username

	type Request struct {
		Username    string `json:"tg_user_name"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	requestData := Request{
		Username:    username,
		Title:       habit.Title,
		Description: habit.Description,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		a.log.Error(fmt.Sprintf("%s: failed to encode to JSON", op), sl.Err(err))
		return 0
	}

	// Send a POST request
	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		a.log.Error(fmt.Sprintf("%s: failed to send http.Post request", op), sl.Err(err))
		return 0
	}
	defer resp.Body.Close()

	response, err := a.readResponse(resp)
	if err != nil {
		a.log.Error(fmt.Sprintf("%s: failed to get the response", op), sl.Err(err))
		return 0
	}

	// Checking if the "habitId" field exists in the response
	habitID, ok := response["habitId"].(float64)
	if !ok {
		a.log.Error(fmt.Sprintf("%s: habitId not found in response", op))
		return 0
	}

	// Converting the float64 habitID to an integer
	habitIDInt := int(habitID)

	a.log.Info(
		fmt.Sprintf("%s: habit created:", op),
		slog.Int("habitId", habitIDInt),
	)

	// // the line bellow only for debugging
	// a.log.Info(fmt.Sprintf("%s: response body", op), slog.Any("value", responseBody))

	return habitIDInt
}

// type getAllHabitsResponse struct {
// 	Data []models.Habit `json:"data"`
// }

// func (h *Handler) getAllHabits(c *gin.Context) {
// 	userId, err := getUserId(c)
// 	if err != nil {
// 		return
// 	}

// 	habits, err := h.services.Habit.GetAll(userId)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getAllHabits: %v", err.Error()))
// 		return
// 	}

// 	c.JSON(http.StatusOK, getAllHabitsResponse{
// 		Data: habits,
// 	})
// }

// func (h *Handler) getHabitById(c *gin.Context) {
// 	userId, err := getUserId(c)
// 	if err != nil {
// 		return
// 	}

// 	habitId, err := getHabitId(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("handler:getHabitById: invalid id param: %v", habitId))
// 		return
// 	}

// 	habit, err := h.services.Habit.GetById(userId, habitId)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getHabitById: %v", err.Error()))
// 		return
// 	}

// 	c.JSON(http.StatusOK, habit)
// }

// func (h *Handler) deleteHabit(c *gin.Context) {
// 	userId, err := getUserId(c)
// 	if err != nil {
// 		return
// 	}

// 	habitId, err := getHabitId(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("handler:deleteHabit: invalid id param: %v", habitId))
// 		return
// 	}

// 	err = h.services.Habit.Delete(userId, habitId)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: deleteHabit: %v", err.Error()))
// 		return
// 	}

// 	c.JSON(http.StatusOK, statusResponse{
// 		Status: "ok",
// 	})
// }

// func (h *Handler) updateHabit(c *gin.Context) {
// 	userId, err := getUserId(c)
// 	if err != nil {
// 		return
// 	}

// 	habitId, err := getHabitId(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("handler:updateHabit: invalid id param: %v", habitId))
// 		return
// 	}

// 	var input models.UpdateHabitInput

// 	if err := c.BindJSON(&input); err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	if err := h.services.Habit.Update(userId, habitId, input); err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: updateHabit: %v", err.Error()))
// 		return
// 	}

// 	c.JSON(http.StatusOK, statusResponse{"ok"})
// }
