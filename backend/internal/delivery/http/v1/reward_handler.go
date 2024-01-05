package v1

import (
	"fmt"
	"net/http"

	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/gin-gonic/gin"
)

// get all rewards of a certain user for a specific reward
func (h *Handler) getPersonalRewardsByHabitId(c *gin.Context) {
	const op = "repository.postgres.reward_handler.getPersonalRewardsByHabitId"

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

	rewards, err := h.services.Reward.GetPersonalRewardsByHabitId(userId, habitId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to get personal rewards by habit id: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get personal rewards by habit id", op), sl.Err(err))
		return
	}

	c.JSON(http.StatusOK, rewards)
}

// get all rewards of a certain user
func (h *Handler) getAllPersonalRewards(c *gin.Context) {
	const op = "repository.postgres.reward_handler.getAllPersonalRewards"

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to get user Id: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get user Id", op), sl.Err(err))
		return
	}

	rewards, err := h.services.Reward.GetAllPersonalRewards(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to get all personal rewards: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get all personal rewards", op), sl.Err(err))
		return
	}

	c.JSON(http.StatusOK, rewards)
}
