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
			habits.GET("/:habitId", h.getHabitById)
			habits.PUT("/:habitId", h.updateHabit)
			habits.DELETE("/:habitId", h.deleteHabit)

			tracker := habits.Group(":habitId/tracker")
			{
				tracker.GET("/", h.getHabitTrackerById)
				tracker.PUT("/", h.updateHabitTracker)
			}

			rewards := habits.Group("/rewards")
			{
				rewards.GET("/", h.getAllPersonalRewards)
				rewards.GET("/:rewardId", h.getPersonalRewardById)
			}
		}

		trackers := api.Group("/trackers")
		{
			trackers.GET("/", h.getAllHabitTrackers)
		}

		admin := api.Group("/admin", h.adminPass)
		{
			users := admin.Group("/users")
			{
				users.GET("/", h.getAllUsers)
				// users.GET("/:id", h.getUserById)
				// users.DELETE("/:id", h.deleteUserById)

				userApi := users.Group("/:userId", h.adminUserPass)
				{
					habits := userApi.Group("/habits")
					{
						habits.POST("/", h.createHabit)
						habits.GET("/", h.getAllHabits)
						habits.GET("/:habitIdAdmin", h.getHabitById)
						habits.PUT("/:habitIdAdmin", h.updateHabit)
						habits.DELETE("/:habitIdAdmin", h.deleteHabit)

						tracker := habits.Group(":habitIdAdmin/tracker")
						{
							tracker.GET("/", h.getHabitTrackerById)
							tracker.PUT("/", h.updateHabitTracker)
						}

						rewards := habits.Group("/rewards")
						{
							rewards.GET("/", h.getAllPersonalRewards)
							rewards.GET("/:rewardIdAdmin", h.getPersonalRewardById)
						}
					}

					trackers := userApi.Group("/trackers")
					{
						trackers.GET("/", h.getAllHabitTrackers)
					}

				}

			}
		}

	}

	return router
}
