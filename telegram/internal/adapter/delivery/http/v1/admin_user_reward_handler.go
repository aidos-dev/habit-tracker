package v1

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/aidos-dev/habit-tracker/backend/internal/models"
// 	"github.com/gin-gonic/gin"
// )

// func (h *Handler) assignReward(c *gin.Context) {
// 	userId, err := getUserId(c)
// 	if err != nil {
// 		return
// 	}

// 	habitId, err := getHabitId(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
// 		return
// 	}

// 	rewardId, err := getRewardId(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
// 		return
// 	}

// 	id, err := h.services.AdminUserReward.AssignReward(userId, habitId, rewardId)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: assignReward: calling from service: %v", err.Error()))
// 		return
// 	}

// 	c.JSON(http.StatusOK, map[string]interface{}{
// 		"id": id,
// 	})
// }

// // take away a reward from a certain user
// func (h *Handler) removeRewardFromUser(c *gin.Context) {
// 	userId, err := getUserId(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: removeRewardFromUser: %v", err.Error()))
// 		return
// 	}

// 	habitId, err := getHabitId(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
// 		return
// 	}

// 	rewardId, err := getRewardId(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
// 		return
// 	}

// 	err = h.services.AdminUserReward.RemoveFromUser(userId, habitId, rewardId)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: removeRewardFromUser: %v", err.Error()))
// 		return
// 	}

// 	c.JSON(http.StatusOK, statusResponse{"ok"})
// }

// // update (replace) a reward of a certain user
// func (h *Handler) updateUserReward(c *gin.Context) {
// 	userId, err := getUserId(c)
// 	if err != nil {
// 		return
// 	}

// 	habitId, err := getHabitId(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
// 		return
// 	}

// 	rewardId, err := getRewardId(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
// 		return
// 	}

// 	var input models.UpdateUserRewardInput
// 	if err := c.BindJSON(&input); err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error from handler: updateUserReward: %v", err.Error()))
// 		return
// 	}

// 	if err := h.services.AdminUserReward.UpdateUserReward(userId, habitId, rewardId, input); err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: updateUserReward: %v", err.Error()))
// 		return
// 	}

// 	c.JSON(http.StatusOK, statusResponse{"ok"})
// }
