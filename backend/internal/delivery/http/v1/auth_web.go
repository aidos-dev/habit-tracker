package v1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUpWeb(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Parsed JSON content: %v\n", input)

	// input = webUserFormat(c, input)

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

func (h *Handler) signInWeb(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("Parsed JSON content: %v\n", input)

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
