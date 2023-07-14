package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aidos-dev/habit-tracker/backend/internal/models"
	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/gin-gonic/gin"
)

/*
adminPass ensures that only admin users have access to admin functionality
*/
func (h *Handler) adminPass(c *gin.Context) {
	const op = "delivery.http.v1.admin_middleware.adminPass"

	userRole, err := getUserRole(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user role not found")
		h.log.Error(fmt.Sprintf("%s: failed to get user role", op), sl.Err(err))
		return
	}

	if userRole != models.Administrator {
		newErrorResponse(c, http.StatusBadRequest, "error: not an admin user: access denied")
		h.log.Error(fmt.Sprintf("%s: not an admin user: access denied", op), sl.Err(err))
		return
	}
}

/*
adminUserPass allows administrator to pass any user id to context
so it can have access and manage users accounts
*/

func (h *Handler) adminUserPass(c *gin.Context) {
	const op = "delivery.http.v1.admin_middleware.adminUserPass"

	userId, err := strconv.ParseFloat(strings.TrimSpace(c.Param("userId")), 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("adminUserIdentity: invalid id param: %v", userId))
		h.log.Error(fmt.Sprintf("%s: invalid id param", op), sl.Err(err))
		return
	}

	c.Set(userCtx, userId)
}
