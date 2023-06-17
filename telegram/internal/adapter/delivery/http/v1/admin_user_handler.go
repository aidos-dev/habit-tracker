package v1

import (
	"fmt"
	"net/http"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/gin-gonic/gin"
)

type getAllUsersResponse struct {
	Data []models.GetUser `json:"data"`
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

func (h *Handler) getUserById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	user, err := h.services.AdminUser.GetUserById(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getUserById: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, user)
}
