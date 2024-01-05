package v1

import (
	"fmt"
	"net/http"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

func (h *Handler) createHabit(c *gin.Context) {
	const op = "delivery.http.v1.habit_handler.createHabit"

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to get user Id: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get user Id", op), sl.Err(err))
		return
	}

	var input models.Habit
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error: failed to get JSON object: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get JSON object", op), sl.Err(err))
		return
	}

	habitId, err := h.services.Habit.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to create a habit: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to create a habit", op), sl.Err(err))
		return
	}

	h.log.Info(
		fmt.Sprintf("%s: habit created:", op),
		slog.Int("habitId", habitId),
		slog.String("input", input.Title),
		slog.String("descpription", input.Description),
		slog.Int("userId", userId),
	)

	c.JSON(http.StatusOK, map[string]interface{}{
		"habitId": habitId,
	})
}

type getAllHabitsResponse struct {
	Data []models.Habit `json:"data"`
}

func (h *Handler) getAllHabits(c *gin.Context) {
	const op = "delivery.http.v1.habit_handler.getAllHabits"

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to get user Id: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get user Id", op), sl.Err(err))
		return
	}

	habits, err := h.services.Habit.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to get habits: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get habits", op), sl.Err(err))
		return
	}

	c.JSON(http.StatusOK, getAllHabitsResponse{
		Data: habits,
	})
}

func (h *Handler) getHabitById(c *gin.Context) {
	const op = "delivery.http.v1.habit_handler.getHabitById"

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

	habit, err := h.services.Habit.GetById(userId, habitId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: habit not found: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to find a habit by Id", op), sl.Err(err))
		return
	}

	c.JSON(http.StatusOK, habit)
}

func (h *Handler) deleteHabit(c *gin.Context) {
	const op = "delivery.http.v1.habit_handler.deleteHabit"

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

	err = h.services.Habit.Delete(userId, habitId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to delete a habit %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to delete a habit", op), sl.Err(err))
		return
	}

	h.log.Info(fmt.Sprintf("%s: a habit is deleted", op), slog.Int("id", habitId))

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) updateHabit(c *gin.Context) {
	const op = "delivery.http.v1.habit_handler.updateHabit"

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

	var input models.UpdateHabitInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error: failed to get JSON object: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get JSON object", op), sl.Err(err))
		return
	}

	if err := h.services.Habit.Update(userId, habitId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to update a habit %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to update a habit", op), sl.Err(err))
		return
	}

	h.log.Info(fmt.Sprintf("%s: a habit has been updated", op), slog.Int("id", habitId))

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
