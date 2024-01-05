package v1

import (
	"fmt"
	"net/http"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

func (h *Handler) assignRole(c *gin.Context) {
	const op = "delivery.http.v1.admin_role_handler.assignRole"

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to get user Id: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get user Id", op), sl.Err(err))
		return
	}

	var input models.UpdateRoleInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error: failed to get JSON object: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get JSON object", op), sl.Err(err))
		return
	}

	id, err := h.services.AdminRole.AssignRole(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to assign role: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to assign role", op), sl.Err(err))
		return
	}

	h.log.Info(
		fmt.Sprintf("%s: a new role has been assigned to user", op),
		slog.Int("user id", id),
		slog.String("role", *input.Role),
	)

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": statusResponse{"ok"},
		"id":     id,
	})
}
