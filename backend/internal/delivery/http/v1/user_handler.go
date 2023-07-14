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
	const op = "delivery.http.v1.user_handler.getUserById"

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to get user Id: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get user Id", op), sl.Err(err))
		return
	}

	user, err := h.services.User.GetUserById(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: a user not found: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to find a user by Id", op), sl.Err(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

type getAllUsersResponse struct {
	Data []models.GetUser `json:"data"`
}

func (h *Handler) getAllUsers(c *gin.Context) {
	const op = "delivery.http.v1.user_handler.getAllUsers"

	users, err := h.services.User.GetAllUsers()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to get all users: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get all users", op), sl.Err(err))
		return
	}

	c.JSON(http.StatusOK, getAllUsersResponse{
		Data: users,
	})
}

func (h *Handler) deleteUser(c *gin.Context) {
	const op = "delivery.http.v1.user_handler.deleteUser"

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: failed to get user Id: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get user Id", op), sl.Err(err))
		return
	}

	deletedUserId, err := h.services.User.DeleteUser(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("%s: failed to delete a user by id: %d: %s", op, userId, err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to delete a user by id: %d", op, userId), sl.Err(err))
		return
	}

	h.log.Info(fmt.Sprintf("%s: user is deleted", op), slog.Int("user id", deletedUserId))

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
