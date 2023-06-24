package v1

import (
	"log"
	"net/http"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUpTelegram(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Parsed JSON content: %v\n", input)

	// input = telegramUserFormat(c, input)

	log.Printf("The TG prepared user is: %v\n", input)

	id, err := h.services.User.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

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
func telegramUserFormat(c *gin.Context, user models.User) models.User {
	var emptyUser models.User // emptyUser created only to return it in case of error

	if user.TgUsername == "" {
		newErrorResponse(c, http.StatusBadRequest, "error: user name is not specified")
		return emptyUser
	}

	if user.Username == "" {
		user.Username = models.Empty
	}

	if user.FirstName == "" {
		user.FirstName = models.Empty
	}

	if user.LastName == "" {
		user.LastName = models.Empty
	}

	if user.Email == "" {
		user.Email = models.Empty
	}

	return user
}

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
