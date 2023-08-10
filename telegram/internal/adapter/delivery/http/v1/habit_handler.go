package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/aidos-dev/habit-tracker/telegram/internal/models"
	"golang.org/x/exp/slog"
)

func (a *AdapterHandler) CreateHabit(username string, habit models.Habit) int {
	const op = "telegram/internal/adapter/delivery/http/v1/habit_handler.CreateHabit"

	a.log.Info(fmt.Sprintf("%s: CreateHabit method called", op))

	// Perform the necessary logic for command1
	a.log.Info(fmt.Sprintf("%s: Executing CreateHabit with text: %s", op, username))

	// Make an HTTP request to the backend service
	requestURL := backendURL + habitsUrl + userQuery + username

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

func (a *AdapterHandler) GetAllHabits(username string) string {
	const op = "telegram/internal/adapter/delivery/http/v1/habit_handler.getAllHabits"

	a.log.Info(fmt.Sprintf("%s: getAllHabits method called", op))

	// Perform the necessary logic for command1
	a.log.Info(fmt.Sprintf("%s: Executing getAllHabits with text: %s", op, username))

	// Make an HTTP request to the backend service
	requestURL := backendURL + habitsUrl + userQuery + username

	// Send a GET request
	resp, err := http.Get(requestURL)
	if err != nil {
		a.log.Error(fmt.Sprintf("%s: failed to send http.Get request", op), sl.Err(err))
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.log.Error(fmt.Sprintf("%s: request failed. status: %d", op, resp.StatusCode), sl.Err(err))
		return ""
	}

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		a.log.Error(fmt.Sprintf("%s: failed to read the response body", op), sl.Err(err))
		return ""
	}

	type allHabits struct {
		Data []models.Habit
	}

	// Decode the response body into the allHabits struct
	var allHabitsData allHabits
	if err := json.Unmarshal(responseBody, &allHabitsData); err != nil {
		a.log.Error(fmt.Sprintf("%s: failed to decode the response", op), sl.Err(err))
		return ""
	}

	a.log.Info(
		fmt.Sprintf("%s: successfully got all habits from backend:", op),
		slog.String("username", username),
		slog.Any("All habits", allHabitsData.Data),
	)

	habitsString := allHabitsToString(allHabitsData.Data)

	// // the line bellow only for debugging
	// a.log.Info(fmt.Sprintf("%s: response body", op), slog.Any("value", responseBody))

	// Return all habits in one string
	return habitsString
}

/*
allHabitsToString converts a slice of all Habits to a nice formatted
list of all habits and collects them into one
string variable to printed out for a telegram user
*/
func allHabitsToString(habitsSlice []models.Habit) string {
	const (
		id      = "Id: "
		habit   = "Habit: "
		desript = "Description: "
		newLine = "\n"
	)

	allHabitsString := ""

	for _, el := range habitsSlice {
		allHabitsString += id
		allHabitsString += strconv.Itoa(el.Id)
		allHabitsString += newLine

		allHabitsString += habit
		allHabitsString += el.Title
		allHabitsString += newLine

		allHabitsString += desript
		allHabitsString += el.Description
		allHabitsString += newLine
		allHabitsString += newLine
	}

	return allHabitsString
}

func (a *AdapterHandler) GetHabitById(habitId int, username string) (models.Habit, error) {
	const op = "telegram/internal/adapter/delivery/http/v1/habit_handler.getHabitById"

	a.log.Info(fmt.Sprintf("%s: getHabitById method called", op))

	// Perform the necessary logic for command1
	a.log.Info(fmt.Sprintf("%s: Executing getHabitById with text: %s", op, username))

	// Make an HTTP request to the backend service
	requestURL := backendURL + habitsUrl + "/" + strconv.Itoa(habitId) + userQuery + username

	var emptyHabit models.Habit

	// Send a GET request
	resp, err := http.Get(requestURL)
	if err != nil {
		a.log.Error(fmt.Sprintf("%s: failed to send http.Get request", op), sl.Err(err))
		return emptyHabit, fmt.Errorf("failed to get a habit")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.log.Error(fmt.Sprintf("%s: request failed. status: %d", op, resp.StatusCode), sl.Err(err))
		return emptyHabit, fmt.Errorf("failed to get a habit")
	}

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		a.log.Error(fmt.Sprintf("%s: failed to read the response body", op), sl.Err(err))
		return emptyHabit, fmt.Errorf("failed to get a habit")
	}

	var habit models.Habit

	// Decode the response body into the habit
	if err := json.Unmarshal(responseBody, &habit); err != nil {
		a.log.Error(fmt.Sprintf("%s: failed to decode the response", op), sl.Err(err))
		return emptyHabit, fmt.Errorf("failed to get a habit")
	}

	a.log.Info(
		fmt.Sprintf("%s: successfully got a habit from backend:", op),
		slog.String("username", username),
		slog.Any("habit", habit),
	)

	// // the line bellow only for debugging
	// a.log.Info(fmt.Sprintf("%s: response body", op), slog.Any("value", responseBody))

	// Return a habit
	return habit, nil
}

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
