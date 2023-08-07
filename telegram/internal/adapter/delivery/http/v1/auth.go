package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aidos-dev/habit-tracker/pkg/loggs/sl"
	"golang.org/x/exp/slog"
)

func (a *AdapterHandler) SignUp(username string) {
	const (
		op        = "telegram/internal/adapter/delivery/http/v1/auth.SignUp"
		signUpUrl = "/auth/sign-up"
	)

	a.log.Info(fmt.Sprintf("%s: SignUp method called", op))

	// Perform the necessary logic for command1
	a.log.Info(fmt.Sprintf("%s: Executing SignUp with text: %s", op, username))

	// Make an HTTP request to the backend service
	requestURL := habitsUrl + signUpUrl

	type Request struct {
		Name string `json:"tg_user_name"`
	}

	requestData := Request{Name: username}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		a.log.Error(fmt.Sprintf("%s: failed to encode to JSON", op), sl.Err(err))
		return
	}

	// Send a POST request
	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		a.log.Error(fmt.Sprintf("%s: failed to send http.Post request", op), sl.Err(err))
		return
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		a.log.Error(fmt.Sprintf("%s: failed to read the responseBody", op), sl.Err(err))
		return
	}

	// the line bellow only for debugging
	a.log.Info(fmt.Sprintf("%s: response body", op), slog.Any("value", responseBody))

	// c.JSON(http.StatusOK, map[string]interface{}{
	// 	"tg_user_name": username,
	// })
}

// func (h *Handler) deleteUser(c *gin.Context) {
// 	userId, err := getUserId(c)
// 	if err != nil {
// 		return
// 	}

// 	deletedUserId, err := h.services.User.DeleteUser(userId)
// 	if err != nil {
// 		newErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("error from handler: delete user: %v", err.Error()))
// 		return
// 	}

// 	response := map[string]any{
// 		"Status":         "ok",
// 		"deleted userId": deletedUserId,
// 	}

// 	c.JSON(http.StatusOK, response)
// }
