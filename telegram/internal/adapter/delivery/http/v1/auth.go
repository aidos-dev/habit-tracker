package v1

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func (a *AdapterHandler) SignUp(username string) {
	log.Print("adapter: SignUp method called")

	// Perform the necessary logic for command1
	log.Println("adapter: SignUp: Executing SignUp with text:", username)

	// Make an HTTP request to the backend service
	requestURL := a.BackendUrl + "/auth/sign-up"

	type Request struct {
		Name string `json:"tg_user_name"`
	}

	requestData := Request{Name: username}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		// c.String(http.StatusInternalServerError, err.Error())
		log.Printf("error: adapter: SignUp: failed to send request: %v", err.Error())
		return
	}

	// Send a POST request
	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		// c.String(http.StatusInternalServerError, err.Error())
		log.Printf("error: %v", err.Error())
		return
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// c.String(http.StatusInternalServerError, err.Error())
		log.Printf("error: %v", err.Error())
		return
	}

	log.Printf("adapter: SignUp: response body: %v", string(responseBody))

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
