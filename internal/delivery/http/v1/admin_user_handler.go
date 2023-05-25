package v1

import (
	"fmt"
	"net/http"

	"github.com/aidos-dev/habit-tracker/internal/models"
	"github.com/gin-gonic/gin"
)

type getAllUsersResponse struct {
	Data []models.User `json:"data"`
}

func (h *Handler) getAllUsers(c *gin.Context) {
	users, err := h.services.AdminUser.GetAllUsers()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getAllUsers: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, getAllUsersResponse{
		Data: users,
	})
}
