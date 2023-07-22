package v1

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	roleCtx             = "userRole"
)

func (h *Handler) userIdentity(c *gin.Context) {
	client := c.Param("client")

	switch client {
	case models.WebClient:
		h.webUserIdentity(c)
	case models.TelegramClient:
		h.tgUserIdentity(c)
	}
}

func (h *Handler) tgUserIdentity(c *gin.Context) {
	const op = "delivery.http.v1.middleware.tgUserIdentity"

	var TgUserName models.TgUser

	if err := c.BindJSON(&TgUserName); err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error: failed to get JSON object: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to get JSON object", op), sl.Err(err))
		return
	}

	user, err := h.services.Authorization.FindTgUser(TgUserName.TgUsername)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: a telegram user not found: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to find a telegram user by tg username", op), sl.Err(err))
		return
	}

	c.Set(userCtx, user.Id)
	c.Set(roleCtx, user.Role)

	c.JSON(http.StatusOK, user)
}

func (h *Handler) webUserIdentity(c *gin.Context) {
	const op = "delivery.http.v1.middleware.webUserIdentity"

	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		h.log.Error(fmt.Sprintf("%s: error", op), "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		h.log.Error(fmt.Sprintf("%s: error", op), "invalid auth header")
		return
	}

	userId, userRole, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		h.log.Error(fmt.Sprintf("%s: failed to parse jwt token", op), sl.Err(err))
		return
	}

	c.Set(roleCtx, userRole)
	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	const op = "delivery.http.v1.middleware.getUserId"

	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found: doesn't exist")
		return 0, fmt.Errorf("%s: user id not found", op)
	}

	idInt := int(id.(float64)) // converting float64 to int

	var n int

	intValue := reflect.TypeOf(idInt) == reflect.TypeOf(n)

	// _, ok = id.(int) // checking if conversion to int was successful
	if !intValue {
		newErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return 0, fmt.Errorf("%s: user id is of invalid type", op)
	}

	return idInt, nil
}

func getUserRole(c *gin.Context) (string, error) {
	const op = "delivery.http.v1.middleware.getUserRole"

	role, exists := c.Get(roleCtx)
	if !exists {
		newErrorResponse(c, http.StatusInternalServerError, "user role not found")
		return "", fmt.Errorf("%s: user role not found", op)
	}

	userRole, ok := role.(string)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user role is of invalid type")
		return "", fmt.Errorf("%s: user role is of invalid type", op)
	}

	return userRole, nil
}

/*
getHabitId returns a habit id depending on the role of a user
(sipmple user or admin). This function is required since parsing
user id and habit id is confusing for c.Param function
*/
func getHabitId(c *gin.Context) (int, error) {
	const op = "delivery.http.v1.middleware.getHabitId"

	userRole, err := getUserRole(c)
	if err != nil {
		return 0, fmt.Errorf("%s: user role not found: %w", op, err)
	}

	var habitId int

	switch userRole {
	case models.UserGeneral:
		habitId, err = strconv.Atoi(c.Param("habitId"))
	case models.Administrator:
		habitId, err = strconv.Atoi(c.Param("habitIdAdmin"))
	}

	return habitId, err
}

/*
getRewardId returns a reward id depending on the role of a user
(sipmple user or admin). This function is required since parsing
user id and reward id is confusing for c.Param function
*/
func getRewardId(c *gin.Context) (int, error) {
	const op = "delivery.http.v1.middleware.getRewardId"

	userRole, err := getUserRole(c)
	if err != nil {
		return 0, fmt.Errorf("%s: user role not found: %w", op, err)
	}

	var rewardId int

	switch userRole {
	case models.UserGeneral:
		rewardId, err = strconv.Atoi(c.Param("rewardId"))
	case models.Administrator:
		rewardId, err = strconv.Atoi(c.Param("rewardIdAdmin"))
	}

	return rewardId, err
}
