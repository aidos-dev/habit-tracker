package v1

import (
	"fmt"
	"net/http"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

func (h *Handler) getUserById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	user, err := h.services.User.GetUserById(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getUserById: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, user)
}

type getAllUsersResponse struct {
	Data []models.GetUser `json:"data"`
}

func (h *Handler) getAllUsers(c *gin.Context) {
	users, err := h.services.User.GetAllUsers()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: getAllUsers: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, getAllUsersResponse{
		Data: users,
	})
}

func (h *Handler) deleteUser(c *gin.Context) {
	const op = "delivery.http.v1.deleteUser"

	userId, err := getUserId(c)
	if err != nil {
		h.log.Error(fmt.Sprintf("%s:failed to find a user by id: %d", op, userId), sl.Err(err))
		return
	}

	deletedUserId, err := h.services.User.DeleteUser(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("%s:failed to delete a user: %d: %s", op, userId, err.Error()))
		h.log.Error(fmt.Sprintf("%s:failed to delete a user: %d", op, userId), sl.Err(err))
		return
	}

	h.log.Info(fmt.Sprintf("%s:user deleted\n", op), slog.Int("id", deletedUserId))

	// c.JSON(http.StatusOK, map[string]interface{}{
	// 	"Status": statusResponse{
	// 		Status: "ok",
	// 	},
	// 	"deletedUserId": deletedUserId,
	// })

	response := map[string]any{
		"Status":         "ok",
		"deleted userId": deletedUserId,
	}

	c.JSON(http.StatusOK, response)
}
