package v1

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/gin-gonic/gin"

	"golang.org/x/exp/slog"
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

	// var TgUserName models.TgUser

	// if err := c.BindJSON(&TgUserName); err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("error: failed to get JSON object: %v", err.Error()))
	// 	h.log.Error(fmt.Sprintf("%s: failed to get JSON object", op), sl.Err(err))
	// 	return
	// }

	TgUserName := c.Query("tgUser")

	h.log.Info(
		fmt.Sprintf("%s: preparing to find a tgUser", op),
		// slog.String("tgUserName", TgUserName.TgUsername),
		slog.String("tgUserName", TgUserName),
	)

	// user, err := h.services.Authorization.FindTgUser(TgUserName.TgUsername)
	user, err := h.services.Authorization.FindTgUser(TgUserName)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error: a telegram user not found: %v", err.Error()))
		h.log.Error(fmt.Sprintf("%s: failed to find a telegram user by tg username", op), sl.Err(err))
		return
	}

	h.log.Info(
		fmt.Sprintf("%s: a tgUser is found and ready to be sent to context", op),
		slog.Any("tgUse", user),
	)

	// h.log.Info(
	// 	fmt.Sprintf("%s: Data type of user.Id", op),
	// 	slog.Any("type", reflect.TypeOf(user.Id)),
	// )

	log.Printf("%s: Data type of user.Id: %v\n", op, reflect.TypeOf(user.Id))

	c.Set(userCtx, user.Id)
	c.Set(roleCtx, user.Role)
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
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header: auth fields are missing")
		h.log.Error(fmt.Sprintf("%s: error", op), "invalid auth header: auth fields are missing")
		return
	}

	if headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header: wrong auth method")
		h.log.Error(fmt.Sprintf("%s: error", op), "invalid auth header: wrong auth method")
		return
	}

	if headerParts[1] == "" {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header: token is missing")
		h.log.Error(fmt.Sprintf("%s: error", op), "invalid auth header: token is missing")
		return
	}

	userId, userRole, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		h.log.Error(fmt.Sprintf("%s: failed to parse jwt token", op), sl.Err(err))
		return
	}

	// h.log.Info(
	// 	fmt.Sprintf("%s: Data type of userId", op),
	// 	slog.Any("type", reflect.TypeOf(userId)),
	// )

	log.Printf("%s: Data type of userId: %v\n", op, reflect.TypeOf(userId))

	c.Set(userCtx, userId)
	c.Set(roleCtx, userRole)
}

func getUserId(c *gin.Context) (int, error) {
	const op = "delivery.http.v1.middleware.getUserId"

	id, ok := c.Get(userCtx)
	if !ok {
		// newErrorResponse(c, http.StatusInternalServerError, "user id not found: doesn't exist")
		return 0, fmt.Errorf("%s: user id not found", op)
	}

	log.Printf("%s: a user id before converting: %d", op, id)

	idInt, err := convertToInt(id)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to convert to int", op)
	}

	return idInt, nil
}

func convertToInt(id any) (int, error) {
	const op = "delivery.http.v1.middleware.convertToInt"
	// idInt := int(id.(float64)) // converting float64 to int
	// idInt, ok := id.(int) // converting  to int
	// if !ok {
	// 	return 0, fmt.Errorf("%s: failed to convert id to int", op)
	// }

	var idInt int

	switch id.(type) {
	case int:
		idInt = id.(int) // converting  to int
		log.Printf("%s: idInt got value from int type: %d", op, id)
	case float64:
		idInt = int(id.(float64)) // converting float64 to int
		log.Printf("%s: idInt got value from float64 type: %d", op, id)
	default:
		return 0, fmt.Errorf("%s: user id is of unknown type", op)
	}

	var n int

	intValue := reflect.TypeOf(idInt) == reflect.TypeOf(n)

	// _, ok = id.(int) // checking if conversion to int was successful
	if !intValue {
		// newErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
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
