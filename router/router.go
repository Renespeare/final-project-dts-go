package routers

import (
	"final-project-dts-go/controllers"
	"final-project-dts-go/middlewares"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	public := router.Group("/api")

	userRouter := public.Group("/users")
	{
		userRouter.POST("/register", controllers.UserRegister)
		userRouter.POST("/login", controllers.UserLogin)
	}

	public.Use(middlewares.Authentication())
	{
		photoRouter := public.Group("/photos")
		{
			photoRouter.POST("/", controllers.CreatePhoto)
			photoRouter.GET("/", controllers.GetAllPhoto)
			photoRouter.GET("/:photoId", controllers.GetOnePhoto)
			photoRouter.PUT("/:photoId", middlewares.PhotoAuthorization(), controllers.UpdatePhoto)
			photoRouter.DELETE("/:photoId", middlewares.PhotoAuthorization(), controllers.DeletePhoto)
		}

		commentRouter := public.Group("/comments")
		{
			commentRouter.POST("/", controllers.CreateComment)
			commentRouter.GET("/", controllers.GetAllComment)
			commentRouter.GET("/:commentId", controllers.GetOneComment)
			commentRouter.PUT("/:commentId", middlewares.CommentAuthorization(), controllers.UpdateComment)
			commentRouter.DELETE("/:commentId", middlewares.CommentAuthorization(), controllers.DeleteComment)
		}

		socialMediaRouter := public.Group("/social_media")
		{
			socialMediaRouter.POST("/", controllers.CreateSocialMedia)
			socialMediaRouter.GET("/", controllers.GetAllSocialMedia)
			socialMediaRouter.GET("/:socialMediaId", controllers.GetOneSocialMedia)
			socialMediaRouter.PUT("/:socialMediaId", middlewares.SocialMediaAuthorization(), controllers.UpdateSocialMedia)
			socialMediaRouter.DELETE("/:socialMediaId", middlewares.SocialMediaAuthorization(), controllers.DeleteSocialMedia)
		}
	}

	return router
}