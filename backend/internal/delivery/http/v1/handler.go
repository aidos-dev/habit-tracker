package v1

import (
	"github.com/aidos-dev/habit-tracker/backend/internal/service"
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

	routerWeb := router.Group("/web")
	routerTelegram := router.Group("/telegram")

	authWeb := routerWeb.Group("/auth")
	{
		authWeb.POST("/sign-up", h.signUpWeb)
		authWeb.POST("/sign-in", h.signInWeb)
	}

	authTelegram := routerTelegram.Group("/auth")
	{
		authTelegram.POST("/sign-up", h.signUpTelegram)
		// authTelegram.GET("/exist", h.tgUserIdentity)
		// authTelegram.POST("/sign-in", h.signInTelegram)
	}

	api := router.Group("/:client/api", h.userIdentity)
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

			rewardsUser := habits.Group(":habitId/rewardsUser")
			{
				rewardsUser.GET("/", h.getPersonalRewardsByHabitId)
			}
		}

		trackers := api.Group("/trackers")
		{
			trackers.GET("/", h.getAllHabitTrackers)
		}

		rewardsUserAll := api.Group("/rewardsUserAll")
		{
			rewardsUserAll.GET("/", h.getAllPersonalRewards)
		}

		userAccount := api.Group("/account")
		{
			userAccount.DELETE("/", h.deleteUser)
		}

		admin := api.Group("/admin", h.adminPass)
		{
			users := admin.Group("/users")
			{
				users.GET("/", h.getAllUsers)

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

						rewardsUser := habits.Group(":habitIdAdmin/rewardsUser")
						{
							rewardsUser.GET("/", h.getPersonalRewardsByHabitId)
						}

						rewardsUserAdmin := habits.Group(":habitIdAdmin/rewardsUserAdmin")
						{
							rewardsUserAdmin.POST("/:rewardIdAdmin", h.assignReward)
							rewardsUserAdmin.PUT("/:rewardIdAdmin", h.updateUserReward)
							rewardsUserAdmin.DELETE("/:rewardIdAdmin", h.removeRewardFromUser)
						}

					}

					trackers := userApi.Group("/trackers")
					{
						trackers.GET("/", h.getAllHabitTrackers)
					}

					rewardsUserAll := userApi.Group("/rewardsUserAll")
					{
						rewardsUserAll.GET("/", h.getAllPersonalRewards)
					}

					roles := userApi.Group("/roles")
					{
						roles.PUT("/", h.assignRole)
					}

					userAccount := userApi.Group("/account")
					{
						userAccount.GET("/", h.getUserById)
						userAccount.DELETE("/", h.deleteUser)

					}

				}

			}

			rewardsAdmin := admin.Group("/rewardsAdmin")
			{
				rewardsAdmin.POST("/", h.createReward)
				rewardsAdmin.GET("/", h.getAllRewards)
				rewardsAdmin.GET("/:rewardId", h.getRewardById)
				rewardsAdmin.PUT("/:rewardId", h.updateReward)
				rewardsAdmin.DELETE("/:rewardId", h.deleteReward)
			}
		}

	}

	return router
}
