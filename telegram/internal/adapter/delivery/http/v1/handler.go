package v1

import "github.com/gin-gonic/gin"

type AdapterHandler struct{}

func NewAdapterHandler() *AdapterHandler {
	return &AdapterHandler{}
}

func (a *AdapterHandler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", a.signUp)
		auth.POST("/sign-in", a.signIn)
	}

	api := router.Group("/api", a.userIdentity)
	{
		habits := api.Group("/habits")
		{
			habits.POST("/", a.createHabit)
			habits.GET("/", a.getAllHabits)
			habits.GET("/:habitId", a.getHabitById)
			habits.PUT("/:habitId", a.updateHabit)
			habits.DELETE("/:habitId", a.deleteHabit)

			tracker := habits.Group(":habitId/tracker")
			{
				tracker.GET("/", a.getHabitTrackerById)
				tracker.PUT("/", a.updateHabitTracker)
			}

			rewardsUser := habits.Group(":habitId/rewardsUser")
			{
				rewardsUser.GET("/", a.getPersonalRewardsByHabitId)
			}
		}

		// trackers := api.Group("/trackers")
		// {
		// 	trackers.GET("/", a.getAllHabitTrackers)
		// }

		// rewardsUserAll := api.Group("/rewardsUserAll")
		// {
		// 	rewardsUserAll.GET("/", a.getAllPersonalRewards)
		// }

		// userAccount := api.Group("/account")
		// {
		// 	userAccount.DELETE("/", a.deleteUser)
		// }

		// admin := api.Group("/admin", a.adminPass)
		// {
		// 	users := admin.Group("/users")
		// 	{
		// 		users.GET("/", a.getAllUsers)

		// 		userApi := users.Group("/:userId", a.adminUserPass)
		// 		{
		// 			habits := userApi.Group("/habits")
		// 			{
		// 				habits.POST("/", a.createHabit)
		// 				habits.GET("/", a.getAllHabits)
		// 				habits.GET("/:habitIdAdmin", a.getHabitById)
		// 				habits.PUT("/:habitIdAdmin", a.updateHabit)
		// 				habits.DELETE("/:habitIdAdmin", a.deleteHabit)

		// 				tracker := habits.Group(":habitIdAdmin/tracker")
		// 				{
		// 					tracker.GET("/", a.getHabitTrackerById)
		// 					tracker.PUT("/", a.updateHabitTracker)
		// 				}

		// 				rewardsUser := habits.Group(":habitIdAdmin/rewardsUser")
		// 				{
		// 					rewardsUser.GET("/", a.getPersonalRewardsByHabitId)
		// 				}

		// 				rewardsUserAdmin := habits.Group(":habitIdAdmin/rewardsUserAdmin")
		// 				{
		// 					rewardsUserAdmin.POST("/:rewardIdAdmin", a.assignReward)
		// 					rewardsUserAdmin.PUT("/:rewardIdAdmin", a.updateUserReward)
		// 					rewardsUserAdmin.DELETE("/:rewardIdAdmin", a.removeRewardFromUser)
		// 				}

		// 			}

		// 			trackers := userApi.Group("/trackers")
		// 			{
		// 				trackers.GET("/", a.getAllHabitTrackers)
		// 			}

		// 			rewardsUserAll := userApi.Group("/rewardsUserAll")
		// 			{
		// 				rewardsUserAll.GET("/", a.getAllPersonalRewards)
		// 			}

		// 			roles := userApi.Group("/roles")
		// 			{
		// 				roles.PUT("/", a.assignRole)
		// 			}

		// 			userAccount := userApi.Group("/account")
		// 			{
		// 				userAccount.GET("/", a.getUserById)
		// 				userAccount.DELETE("/", a.deleteUser)

		// 			}

		// 		}

		// 	}

		// 	rewardsAdmin := admin.Group("/rewardsAdmin")
		// 	{
		// 		rewardsAdmin.POST("/", a.createReward)
		// 		rewardsAdmin.GET("/", a.getAllRewards)
		// 		rewardsAdmin.GET("/:rewardId", a.getRewardById)
		// 		rewardsAdmin.PUT("/:rewardId", a.updateReward)
		// 		rewardsAdmin.DELETE("/:rewardId", a.deleteReward)
		// 	}
		// }

	}

	return router
}
