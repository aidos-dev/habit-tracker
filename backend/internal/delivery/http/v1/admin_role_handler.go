package v1

import (
	"fmt"
	"net/http"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) assignRole(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: assignRole: user not found: %v", err.Error()))
		return
	}

	var input models.UpdateRoleInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error from handler: assignRole: %v", err.Error()))
		return
	}

	id, err := h.services.AdminRole.AssignRole(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: updateReward: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": statusResponse{"ok"},
		"id":     id,
	})
}
