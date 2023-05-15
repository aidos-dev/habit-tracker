package v1

import (
	"github.com/aidos-dev/habit-tracker/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		habits := api.Group("/habits")
		{
			habits.POST("/", h.createHabit)
			habits.GET("/", h.getAllHabits)
			habits.GET("/:id", h.getHabitById)
			habits.PUT("/:id", h.updateHabit)
			habits.DELETE("/:id", h.deleteHabit)

			tracker := habits.Group(":id/tracker")
			{
				tracker.GET("/", h.getHabitTrackerById)
				tracker.PUT("/", h.updateHabitTracker)
			}
		}

		// admin := api.Group("/admin", h.adminAuth)
		// {
		// }
	}

	return router
}
