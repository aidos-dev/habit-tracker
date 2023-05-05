package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/aidos-dev/habit-tracker"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createReward(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		return
	}

	var input habit.Reward
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error from handler: createReward: %v", err.Error()))
		return
	}

	id, err := h.services.Reward.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: createReward: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

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

	var input habit.Reward
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error from handler: assignReward: %v", err.Error()))
		return
	}

	id, err := h.services.Reward.AssignReward(userId, rewardId, habitId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: assignReward: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// get all rewards from reward table. Independent objects, not
// associated with users
func (h *Handler) getAllRewards(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getAllRewards: %v", err.Error()))
		return
	}

	rewards, err := h.services.Reward.GetAllRewards()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getAllRewards: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, rewards)
}

// get reward from reward table. Independent object, not
// associated with users
func (h *Handler) getRewardById(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getRewardById: %v", err.Error()))
		return
	}

	rewardId, err := strconv.Atoi(c.Param("rewardId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	reward, err := h.services.Reward.GetById(rewardId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getRewardById: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, reward)
}

// get all rewards of a certain user
func (h *Handler) getRewardsByUserId(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getRewardsByUserId: %v", err.Error()))
		return
	}

	rewards, err := h.services.Reward.GetByUserId(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getRewardsByUserId: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, rewards)
}

// delete reward from reward table. Independent object, not
// associated with users
func (h *Handler) deleteReward(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: deleteReward: %v", err.Error()))
		return
	}

	rewardId, err := strconv.Atoi(c.Param("rewardId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	err = h.services.Reward.Delete(rewardId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: deleteReward: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
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

	err = h.services.Reward.RemoveFromUser(userId, rewardId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: removeRewardFromUser: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// update a reward in reward table. Independent object, not
// associated with users
func (h *Handler) updateReward(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		return
	}

	rewardId, err := strconv.Atoi(c.Param("rewardId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input habit.UpdateRewardInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error from handler: updateReward: %v", err.Error()))
		return
	}

	if err := h.services.Reward.UpdateReward(rewardId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: updateReward: %v", err.Error()))
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

	var input habit.UpdateUserRewardInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error from handler: updateUserReward: %v", err.Error()))
		return
	}

	if err := h.services.Reward.UpdateUserReward(userId, rewardId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: updateUserReward: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
