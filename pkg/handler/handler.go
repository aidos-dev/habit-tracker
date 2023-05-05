package handler

import (
	"github.com/aidos-dev/habit-tracker/pkg/service"
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

			trackers := habits.Group(":id/trackers")
			{
				// trackers.POST("/", h.createHabitTracker) // temporarily disabled
				trackers.GET("/", h.getAllHabitTrackers)
			}
		}

		trackers := api.Group("trackers")
		{
			trackers.GET("/:id", h.getHabitTrackerById)
			trackers.PUT("/:id", h.updateHabitTracker)
			// trackers.DELETE("/:id", h.deleteHabitTracker) // temporarily disabled
		}

		// admin := api.Group("/admin", h.adminAuth)
		// {
		// }
	}

	return router
}
