package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/aidos-dev/habit-tracker"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createHabitTracker(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userHabitId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	var input habit.HabitTracker
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error from handler: createHabitTracker: %v", err.Error()))
		return
	}

	id, err := h.services.HabitTracker.Create(userHabitId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllHabitTrackers(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	trackers, err := h.services.HabitTracker.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getAllHabitTrackers: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, trackers)
}

func (h *Handler) getHabitTrackerById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	habitId, err := strconv.Atoi(c.Param("habitId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	tracker, err := h.services.HabitTracker.GetById(userId, habitId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getHabitTrackerById: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, tracker)
}

func (h *Handler) deleteHabitTracker(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	habitId, err := strconv.Atoi(c.Param("habitId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	err = h.services.HabitTracker.Delete(userId, habitId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: deleteHabitTracker: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) updateHabitTracker(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	habitId, err := strconv.Atoi(c.Param("habitId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input habit.UpdateTrackerInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.HabitTracker.Update(userId, habitId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: updateHabitTracker: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
