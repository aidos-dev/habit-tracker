package v1

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	roleCtx             = "userRole"
)

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

func (h *Handler) adminPass(c *gin.Context) {
}

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
