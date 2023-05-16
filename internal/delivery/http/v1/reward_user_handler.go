package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) assignReward(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	rewardId, err := strconv.Atoi(c.Param("rewardId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	habitId, err := strconv.Atoi(c.Param("habitId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input models.Reward
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error from handler: assignReward: %v", err.Error()))
		return
	}

	id, err := h.services.UserReward.AssignReward(userId, rewardId, habitId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: assignReward: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// get all rewards of a certain user
func (h *Handler) getRewardsByUserId(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getRewardsByUserId: %v", err.Error()))
		return
	}

	rewards, err := h.services.UserReward.GetByUserId(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getRewardsByUserId: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, rewards)
}

// take away a reward from a certain user
func (h *Handler) removeRewardFromUser(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: removeRewardFromUser: %v", err.Error()))
		return
	}

	rewardId, err := strconv.Atoi(c.Param("rewardId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	err = h.services.UserReward.RemoveFromUser(userId, rewardId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: removeRewardFromUser: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// update (replace) a reward of a certain user
func (h *Handler) updateUserReward(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	rewardId, err := strconv.Atoi(c.Param("rewardId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input models.UpdateUserRewardInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error from handler: updateUserReward: %v", err.Error()))
		return
	}

	if err := h.services.UserReward.UpdateUserReward(userId, rewardId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: updateUserReward: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
