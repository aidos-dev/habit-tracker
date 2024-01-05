package v1

import (
	"fmt"
	"net/http"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

func (h *Handler) getAllHabitTrackers(c *gin.Context) {
	const op = "delivery.http.v1.habit_tracker_handler.getAllHabitTrackers"

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to get user Id: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get user Id", op), sl.Err(err))
		return
	}

	trackers, err := h.services.HabitTracker.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to get all habit trackers: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get all habit trackers", op), sl.Err(err))
		return
	}

	c.JSON(http.StatusOK, trackers)
}

func (h *Handler) getHabitTrackerById(c *gin.Context) {
	const op = "delivery.http.v1.habit_tracker_handler.getHabitTrackerById"

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to get user Id: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get user Id", op), sl.Err(err))
		return
	}

	habitId, err := getHabitId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: invalid id param: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get habit Id", op), sl.Err(err))
		return
	}

	tracker, err := h.services.HabitTracker.GetById(userId, habitId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: habit tracker not found: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to find habit tracker by Id", op), sl.Err(err))
		return
	}

	c.JSON(http.StatusOK, tracker)
}

func (h *Handler) updateHabitTracker(c *gin.Context) {
	const op = "delivery.http.v1.habit_tracker_handler.updateHabitTracker"

	h.log.Info(
		fmt.Sprintf("%s: tracker update method was called", op),
	)

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to get user Id: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get user Id", op), sl.Err(err))
		return
	}

	habitId, err := getHabitId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: invalid id param: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get habit Id", op), sl.Err(err))
		return
	}

	var input models.UpdateTrackerInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error: failed to get JSON object: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get JSON object", op), sl.Err(err))
		return
	}

	if err := h.services.HabitTracker.Update(userId, habitId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to update a habit tracker %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to update a habit tracker", op), sl.Err(err))
		return
	}

	h.log.Info(
		fmt.Sprintf("%s: a habit tracker has been updated", op),
		slog.Int("habit id", habitId),
	)

	c.JSON(http.StatusOK, statusResponse{"ok"})

	// c.JSON(http.StatusOK, map[string]interface{}{
	// 	"status": "ok",
	// })
}

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
