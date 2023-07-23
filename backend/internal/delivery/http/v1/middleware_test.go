package v1

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/aidos-dev/habit-tracker/backend/internal/service"
	mock_service "github.com/aidos-dev/habit-tracker/backend/internal/service/mocks"
	"github.com/aidos-dev/habit-tracker/pkg/loggs/handlers/slogdiscard"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func Test_handler_webUserIdentity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string)

	type ResponseBody struct {
		UserID   int    `json:"userId"`
		UserRole string `json:"userRole"`
	}

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody ResponseBody
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, "test_role", nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: ResponseBody{
				UserID:   1,
				UserRole: "test_role",
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.token)

			log := slogdiscard.NewDiscardLogger()

			services := &service.Service{Authorization: auth}
			handler := NewHandler(log, services)

			// Init Endpoint
			r := gin.New()
			r.GET("/identity", handler.webUserIdentity, func(c *gin.Context) {
				// Get user ID from the context
				userID, _ := c.Get(userCtx)

				// Get user role from the context
				userRole, _ := c.Get(roleCtx)

				// Create a response struct with userID and userRole
				responseBody := ResponseBody{
					UserID:   userID.(int),      // Convert userID to int
					UserRole: userRole.(string), // Convert userRole to string
				}

				// Return the response as JSON
				c.JSON(200, responseBody)
			})

			// Init Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/identity", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			// Make Request
			r.ServeHTTP(w, req)

			if w.Code != testCase.expectedStatusCode {
				t.Errorf("Expected status code: %d but got: %d", testCase.expectedStatusCode, w.Code)
				// t.Error("Status codes don’t match")
				// t.Logf("Expected status code: %d", testCase.expectedStatusCode)
				// t.Logf("But got: %d", w.Code)
			}

			var responseBody ResponseBody
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			if err != nil {
				t.Errorf("Failed to unmarshal response body: %v", err)
			}

			if responseBody.UserID != testCase.expectedResponseBody.UserID {
				t.Errorf("Expected userID: %d but got: %d", testCase.expectedResponseBody.UserID, responseBody.UserID)
			}

			if responseBody.UserRole != testCase.expectedResponseBody.UserRole {
				t.Errorf("Expected userRole '%s' but got '%s'", testCase.expectedResponseBody.UserRole, responseBody.UserRole)
			}
		})
	}
}
