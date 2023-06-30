package v1

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aidos-dev/habit-tracker/telegram/internal/models"
)

func (a *AdapterHandler) UpdateHabitTracker(username string, habitId string, habitTracker models.HabitTracker) {
	log.Print("adapter UpdateHabitTracker method called")

	// Perform the necessary logic for command1
	log.Println("Executing UpdateHabitTracker with text:", username)

	// Make an HTTP request to the backend service
	requestURL := a.BackendUrl + "/api/habits/" + habitId + "/tracker"

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

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		// c.String(http.StatusInternalServerError, err.Error())
		log.Printf("error: failed to send request: %v", err.Error())
		return
	}

	// Send a PUT request
	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(requestBody))
	if err != nil {
		// c.String(http.StatusInternalServerError, err.Error())
		log.Printf("error: %v", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error executing request: %v", err)
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
