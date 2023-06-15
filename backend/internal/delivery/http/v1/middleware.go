package v1

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/gin-gonic/gin"
)

const (
	webClient      = "webClient"
	telegramClient = "telegramClient"
	// clientCtx           = "client"
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	roleCtx             = "userRole"
)

// func (h *Handler) clientType(c *gin.Context) {
// 	var client models.Client

// 	if err := c.BindJSON(&client); err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	c.Set(clientCtx, client.ClientType)
// }

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	claims, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	data := claims["data"].(map[string]any)

	userId := data["userId"].(float64)
	userRole := data["userRole"].(string)

	c.Set(roleCtx, userRole)
	c.Set(userCtx, userId)
}

// func getClientType(c *gin.Context) (string, error) {
// 	clienCtx, exists := c.Get(clientCtx)
// 	if !exists {
// 		newErrorResponse(c, http.StatusInternalServerError, "client type not found: doesn't exist")
// 		return "", errors.New("client type not found: doesn't exist")
// 	}

// 	fmt.Printf("the data type of clientCtx: %v\n", reflect.TypeOf(clienCtx))

// 	clientType, stringValue := clienCtx.(string)
// 	fmt.Printf("the data type of clientType: %v\n", reflect.TypeOf(clientType))
// 	if !stringValue {
// 		newErrorResponse(c, http.StatusInternalServerError, "client type is of invalid type")
// 		return "", errors.New("error: client type is of invalid type")
// 	}

// 	return clientType, nil
// }

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found: doesn't exist")
		return 0, errors.New("user id not found: doesn't exist")
	}

	idInt := int(id.(float64)) // converting float64 to int

	var n int

	intValue := reflect.TypeOf(idInt) == reflect.TypeOf(n)

	// _, ok = id.(int) // checking if conversion to int was successful
	if !intValue {
		newErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return 0, errors.New("user id not found: user id is of invalid type")
	}

	return idInt, nil
}

func getUserRole(c *gin.Context) (string, error) {
	role, exists := c.Get(roleCtx)
	if !exists {
		newErrorResponse(c, http.StatusInternalServerError, "user role not found: doesn't exist")
		return "", errors.New("user role not found: : doesn't exist")
	}

	userRole, ok := role.(string)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user role is of invalid type")
		return "", errors.New("user role not found: user role is of invalid type")
	}

	return userRole, nil
}

/*
getHabitId returns a habit id depending on the role of a user
(sipmple user or admin). This function is required since parsing
user id and habit id is confusing for c.Param function
*/
func getHabitId(c *gin.Context) (int, error) {
	userRole, err := getUserRole(c)

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
	userRole, err := getUserRole(c)

	var rewardId int

	switch userRole {
	case models.UserGeneral:
		rewardId, err = strconv.Atoi(c.Param("rewardId"))
	case models.Administrator:
		rewardId, err = strconv.Atoi(c.Param("rewardIdAdmin"))
	}

	return rewardId, err
}
