package v1

import (
	"log"
	"net/http"

	"github.com/aidos-dev/habit-tracker/telegram/internal/models"
	"github.com/gin-gonic/gin"
)

func (a *AdapterHandler) FindTgUser(c *gin.Context, username string, userExists *bool) {
	c.String(http.StatusOK, username)

	var TgUserName models.TgUserName

	if err := c.BindJSON(&TgUserName); err != nil {
		log.Printf("error: FindTgUser: %v\n", err.Error())
		*userExists = false
	}
	log.Printf("Parsed JSON content: %v\n", TgUserName)

	if TgUserName.Username != "" {
		*userExists = true
		return
	}
}
