package v1

import "github.com/gin-gonic/gin"

const backendURL = "http://habit-tracker:8000"

type AdapterHandler struct {
	Engine *gin.Engine
	// Router     *gin.RouterGroup
	BackendUrl string
}

func NewAdapterHandler() *AdapterHandler {
	return &AdapterHandler{
		Engine:     gin.New(),
		BackendUrl: backendURL,
	}
}

// func HandleInternalCommand(router *gin.Engine, backendURL string, command string, method string) {
// 	// Register the handler function for the specific command and method
// 	router.Handle(method, fmt.Sprintf("/internal-command/%s", command), func(c *gin.Context) {
// 		// Create a new request to the backend service
// 		req, err := http.NewRequest(method, backendURL, c.Request.Body)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
// 			return
// 		}

// 		// Forward the request to the backend service
// 		client := http.DefaultClient
// 		resp, err := client.Do(req)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request to backend"})
// 			return
// 		}
// 		defer resp.Body.Close()

// 		// Read the response from the backend service
// 		body, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
// 			return
// 		}

// 		log.Printf("response body: %v", string(body))

// 		// Set the response from the backend service as the response for the current request
// 		c.Status(resp.StatusCode)
// 		c.Writer.Write(body)
// 	})
// }
