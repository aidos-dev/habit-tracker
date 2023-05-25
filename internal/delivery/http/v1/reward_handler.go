package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// get all rewards of a certain user
func (h *Handler) getPersonalRewardById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getRewardsByUserId: %v", err.Error()))
		return
	}

	rewards, err := h.services.Reward.GetPersonalRewardById(userId, 0)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getRewardsByUserId: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, rewards)
}

func (h *Handler) getAllPersonalRewards(c *gin.Context) {
}
