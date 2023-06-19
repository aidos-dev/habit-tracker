package v1

// import (
// 	"fmt"
// 	"net/http"
// 	"strconv"
// 	"strings"

// 	"github.com/aidos-dev/habit-tracker/backend/internal/models"
// 	"github.com/gin-gonic/gin"
// )

// /*
// adminPass ensures that only admin users have access to admin functionality
// */
// func (h *Handler) adminPass(c *gin.Context) {
// 	userRole, err := getUserRole(c)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, "user role not found: doesn't exist")
// 		return
// 	}

// 	if userRole != models.Administrator {
// 		newErrorResponse(c, http.StatusBadRequest, "error: not an admin user: access denied")
// 		return
// 	}
// }

// /*
// adminUserPass allows administrator to pass any user id to context
// so it can have access and manage users accounts
// */

// func (h *Handler) adminUserPass(c *gin.Context) {
// 	userId, err := strconv.ParseFloat(strings.TrimSpace(c.Param("userId")), 64)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("handler:adminUserIdentity: invalid id param: %v", userId))
// 		return
// 	}

// 	c.Set(userCtx, userId)
// }
