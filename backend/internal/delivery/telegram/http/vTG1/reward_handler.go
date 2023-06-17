package vTG1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// get all rewards of a certain user for a specific reward
func (h *Handler) getPersonalRewardsByHabitId(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getPersonalRewardsByHabitId: invalid userId param: %v", err.Error()))
		return
	}

	habitId, err := getHabitId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error from handler: getPersonalRewardsByHabitId: invalid habitId param: %v", habitId))
		return
	}

	rewards, err := h.services.Reward.GetPersonalRewardsByHabitId(userId, habitId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getPersonalRewardsByHabitId: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, rewards)
}

// get all rewards of a certain user
func (h *Handler) getAllPersonalRewards(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getAllPersonalRewards: invalid userId param: %v", err.Error()))
		return
	}

	rewards, err := h.services.Reward.GetAllPersonalRewards(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getAllPersonalRewards: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, rewards)
}
