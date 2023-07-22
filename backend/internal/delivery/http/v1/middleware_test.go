package v1

import (
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

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, "test_role", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
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
				id, _ := c.Get(userCtx)
				c.String(200, "%d", id)
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

			if w.Body.String() != testCase.expectedResponseBody {
				t.Errorf("Expected response body '%s' but got '%s'", testCase.expectedResponseBody, w.Body.String())
			}
		})
	}
}
