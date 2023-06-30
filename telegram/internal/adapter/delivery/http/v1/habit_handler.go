package v1

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aidos-dev/habit-tracker/telegram/internal/models"
)

func (a *AdapterHandler) CreateHabit(username string, habit models.Habit) {
	log.Print("adapter CreateHabit method called")

	// Perform the necessary logic for command1
	log.Println("Executing CreateHabit with text:", username)

	// Make an HTTP request to the backend service
	requestURL := a.BackendUrl + "/api/habits"

	type Request struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	requestData := Request{
		Title:       habit.Title,
		Description: habit.Description,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		// c.String(http.StatusInternalServerError, err.Error())
		log.Printf("error: failed to send request: %v", err.Error())
		return
	}

	// Send a POST request
	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		// c.String(http.StatusInternalServerError, err.Error())
		log.Printf("error: %v", err.Error())
		return
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// c.String(http.StatusInternalServerError, err.Error())
		log.Printf("error: %v", err.Error())
		return
	}

	log.Printf("response body: %v", string(responseBody))
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
