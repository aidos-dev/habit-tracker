package v1

import (
	"fmt"
	"net/http"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

func (h *Handler) assignReward(c *gin.Context) {
	const op = "delivery.http.v1.admin_user_reward_handler.assignReward"

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

	rewardId, err := getRewardId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: reward not found: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get reward Id", op), sl.Err(err))
		return
	}

	id, err := h.services.AdminUserReward.AssignReward(userId, habitId, rewardId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to assign reward: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to assign reward", op), sl.Err(err))
		return
	}

	h.log.Info(
		fmt.Sprintf("%s: a new reward has been assigned to user", op),
		slog.Int("user id", userId),
		slog.Int("habit id", habitId),
		slog.Int("reward id", rewardId),
	)

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// take away a reward from a certain user
func (h *Handler) removeRewardFromUser(c *gin.Context) {
	const op = "delivery.http.v1.admin_user_reward_handler.removeRewardFromUser"

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

	rewardId, err := getRewardId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: reward not found: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get reward Id", op), sl.Err(err))
		return
	}

	err = h.services.AdminUserReward.RemoveFromUser(userId, habitId, rewardId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to remove reward from user %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to remove reward from user", op), sl.Err(err))
		return
	}

	h.log.Info(
		fmt.Sprintf("%s: a reward is removed from user", op),
		slog.Int("user id", userId),
		slog.Int("habit id", habitId),
		slog.Int("reward id", rewardId),
	)

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// update (replace) a reward of a certain user
func (h *Handler) updateUserReward(c *gin.Context) {
	const op = "delivery.http.v1.admin_user_reward_handler.updateUserReward"

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

	rewardId, err := getRewardId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: reward not found: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get reward Id", op), sl.Err(err))
		return
	}

	var input models.UpdateUserRewardInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error: failed to get JSON object: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get JSON object", op), sl.Err(err))
		return
	}

	if err := h.services.AdminUserReward.UpdateUserReward(userId, habitId, rewardId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to assign reward: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to assign reward", op), sl.Err(err))
		return
	}

	h.log.Info(
		fmt.Sprintf("%s: a new reward has been assigned to user", op),
		slog.Int("user id", userId),
		slog.Int("habit id", habitId),
		slog.Int("reward id", rewardId),
	)

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
