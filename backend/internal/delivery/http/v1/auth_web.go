package v1

import (
	"fmt"
	"net/http"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

func (h *Handler) signUpWeb(c *gin.Context) {
	const op = "delivery.http.v1.signUpWeb"

	var input models.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error(fmt.Sprintf("%s:failed to get JSON object\n", op), sl.Err(err))
		return
	}

	// the line bellow only for debugging
	// h.log.Info("Parsed JSON content", slog.Any("value", input))

	id, err := h.services.User.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error(fmt.Sprintf("%s: failed to add new user", op), sl.Err(err))
		return
	}

	h.log.Info(fmt.Sprintf("%s:user has been added\n", op), slog.Int("id", id))

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signInWeb(c *gin.Context) {
	const op = "delivery.http.v1.signInWeb"

	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error(fmt.Sprintf("%s:failed to get JSON object\n", op), sl.Err(err))
		return
	}

	// the line bellow only for debugging
	// h.log.Info("Parsed JSON content", slog.Any("value", input))

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		h.log.Error(fmt.Sprintf("%s:failed to generate JWT token", op), sl.Err(err))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

// func (h *Handler) deleteUser(c *gin.Context) {
// 	const op = "delivery.http.v1.deleteUser"

// 	userId, err := getUserId(c)
// 	if err != nil {
// 		h.log.Error(fmt.Sprintf("%s:failed to find a user by id: %d", op, userId), sl.Err(err))
// 		return
// 	}

// 	deletedUserId, err := h.services.User.DeleteUser(userId)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("%s:failed to delete a user: %d: %s", op, userId, err.Error()))
// 		h.log.Error(fmt.Sprintf("%s:failed to delete a user: %d", op, userId), sl.Err(err))
// 		return
// 	}

// 	h.log.Info(fmt.Sprintf("%s:user deleted\n", op), slog.Int("id", deletedUserId))

// 	// c.JSON(http.StatusOK, map[string]interface{}{
// 	// 	"Status": statusResponse{
// 	// 		Status: "ok",
// 	// 	},
// 	// 	"deletedUserId": deletedUserId,
// 	// })

// 	response := map[string]any{
// 		"Status":         "ok",
// 		"deleted userId": deletedUserId,
// 	}

// 	c.JSON(http.StatusOK, response)
// }

/*
webUserFormat prepares user input to be
registered as a web user, using all required credentials. Also
it created a telegram userName as NULL value just as a placeholder. In the future a web user
will be able to replace it with real telegram userName
*/
func webUserFormat(c *gin.Context, user models.User) models.User {
	user.TgUsername = models.Empty

	return user
}
