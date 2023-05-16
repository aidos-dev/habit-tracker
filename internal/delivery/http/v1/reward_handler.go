package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createReward(c *gin.Context) {
	_, err := getUserId(c)
	if err != nil {
		return
	}

	var input models.Reward
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

	var input models.UpdateRewardInput
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
