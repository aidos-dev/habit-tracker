package v1

import (
	"fmt"
	"net/http"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

func (h *Handler) signUpTelegram(c *gin.Context) {
	const op = "delivery.http.v1.auth_telegram.signUpTelegram"

	var input models.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error: failed to get JSON object: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get JSON object", op), sl.Err(err))
		return
	}

	if err := input.Validate(); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error(fmt.Sprintf("%s: invalid input", op), sl.Err(err))
		return
	}

	id, err := h.services.User.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		h.log.Error(fmt.Sprintf("%s: failed to add new user", op), sl.Err(err))
		return
	}

	h.log.Info(fmt.Sprintf("%s: a new user has been added", op), slog.Int("id", id))

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// /*
// tgUser gets telegram user name from url and puts it
// to gin context
// */
// func (h *Handler) tgUser(c *gin.Context) {
// 	tgUserName := strings.TrimSpace(c.Param("tgUser"))

// 	c.Set(tgUserCtx, tgUserName)
// }

// func (h *Handler) signInTelegram(c *gin.Context) {
// 	var input signInInput

// 	if err := c.BindJSON(&input); err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	c.JSON(http.StatusOK, map[string]interface{}{
// 		"token": token,
// 	})
// }

/*
telegramUserFormat prepares user input to be
registered as a telegram user, using only
telegram username and NULL values in other fields. Otherwise the repository layer
will not allow to create a user without other credentials
*/
// func telegramUserFormat(c *gin.Context, user models.User) models.User {
// 	const op = "delivery.http.v1.auth_telegram.telegramUserFormat"

// 	var emptyUser models.User // emptyUser created only to return it in case of error

// 	if user.TgUsername == "" {
// 		newErrorResponse(c, http.StatusBadRequest, "error: user name is not specified")
// 		return emptyUser
// 	}

// 	if user.Username == "" {
// 		user.Username = models.Empty
// 	}

// 	if user.FirstName == "" {
// 		user.FirstName = models.Empty
// 	}

// 	if user.LastName == "" {
// 		user.LastName = models.Empty
// 	}

// 	if user.Email == "" {
// 		user.Email = models.Empty
// 	}

// 	return user
// }

// func (h *Handler) deleteUserTelegram(c *gin.Context) {
// 	userId, err := getUserId(c)
// 	if err != nil {
// 		return
// 	}

// 	deletedUserId, err := h.services.User.DeleteUser(userId)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: delete user: %v", err.Error()))
// 		return
// 	}

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
