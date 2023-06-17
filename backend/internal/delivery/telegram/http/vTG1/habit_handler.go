package vTG1

import (
	"fmt"
	"net/http"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createHabit(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input models.Habit
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	habitId, err := h.services.Habit.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"habitId": habitId,
	})
}

type getAllHabitsResponse struct {
	Data []models.Habit `json:"data"`
}

func (h *Handler) getAllHabits(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	habits, err := h.services.Habit.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getAllHabits: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, getAllHabitsResponse{
		Data: habits,
	})
}

func (h *Handler) getHabitById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	habitId, err := getHabitId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("handler:getHabitById: invalid id param: %v", habitId))
		return
	}

	habit, err := h.services.Habit.GetById(userId, habitId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getHabitById: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, habit)
}

func (h *Handler) deleteHabit(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	habitId, err := getHabitId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("handler:deleteHabit: invalid id param: %v", habitId))
		return
	}

	err = h.services.Habit.Delete(userId, habitId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: deleteHabit: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) updateHabit(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	habitId, err := getHabitId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("handler:updateHabit: invalid id param: %v", habitId))
		return
	}

	var input models.UpdateHabitInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Habit.Update(userId, habitId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: updateHabit: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
