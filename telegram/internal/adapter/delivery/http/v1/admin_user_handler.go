package v1

import (
	"fmt"
	"net/http"

	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"github.com/aidos-dev/habit-tracker/telegram/internal/models"
	"github.com/gin-gonic/gin"
)

func (a *AdapterHandler) FindTgUser(c *gin.Context, username string, userExists *bool) {
	const op = "telegram/internal/adapter/delivery/http/v1/admin_user_handler.FindTgUser"

	var TgUserName models.TgUserName

	if err := c.BindJSON(&TgUserName); err != nil {
		// newErrorResponse(c, http.StatusBadRequest, err.Error())
		a.log.Error(fmt.Sprintf("%s: failed to get JSON object", op), sl.Err(err))
		*userExists = false
	}

	// the line bellow only for debugging
	// a.log.Info(fmt.Sprintf("%s: Parsed JSON content", op), slog.Any("value", TgUserName))

	if TgUserName.Username != "" {
		*userExists = true
		return
	}

	c.String(http.StatusOK, username)
}
