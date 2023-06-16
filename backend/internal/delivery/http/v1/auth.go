package v1

import (
	"fmt"
	"net/http"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var authStruct models.Auth

	if err := c.BindJSON(&authStruct); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Printf("Parsed JSON content: %v\n", authStruct)

	clientType := authStruct.Client.ClientType

	input := authStruct.User

	fmt.Printf("The client type is: %v\n", clientType)

	fmt.Printf("The user name is: %v\n", input)

	input = prepareUserByClient(c, input, clientType)

	fmt.Printf("The TG prepared user is: %v\n", input)

	id, err := h.services.User.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (h *Handler) deleteUser(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	deletedUserId, err := h.services.User.DeleteUser(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: delete user: %v", err.Error()))
		return
	}

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

func prepareUserByClient(c *gin.Context, user models.User, clientType string) models.User {
	switch clientType {
	case webClient:
		user = webUserFormat(c, user)
	case telegramClient:
		user = telegramUserFormat(c, user)
	}

	return user
}

/*
webUserFormat prepares user input to be
registered as a web user, using all required credentials. Also
it created a telegram userName as a copy of web
userName just as a placeholder. In the future a web user
will be able to replace it with real telegram userName
*/
func webUserFormat(c *gin.Context, user models.User) models.User {
	user.TgUsername = fmt.Sprintf("copy_u:%s", user.Username)

	return user
}

/*
telegramUserFormat prepares user input to be
registered as a telegram user, using only
user name. Otherwise the repository layer will not
allow to create a user without other credentials
*/
func telegramUserFormat(c *gin.Context, user models.User) models.User {
	var emptyUser models.User // emptyUser created only to return it in case of error

	if user.TgUsername == "" {
		newErrorResponse(c, http.StatusBadRequest, "error: user name is not specified")
		return emptyUser
	}

	if user.Username == "" {
		user.Username = fmt.Sprintf("copy_tg:%s", user.TgUsername)
	}

	if user.FirstName == "" {
		user.FirstName = fmt.Sprintf("c_tg:%s", user.TgUsername)
	}

	if user.LastName == "" {
		user.LastName = fmt.Sprintf("c_tg:%s", user.TgUsername)
	}

	if user.Email == "" {
		user.Email = fmt.Sprintf("c_tg:%s", user.TgUsername)
	}

	return user
}
