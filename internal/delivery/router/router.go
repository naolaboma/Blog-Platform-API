package router

import (
	"Blog-API/internal/delivery/controllers"
	"Blog-API/internal/infrastructure/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *controllers.UserHandler,
	blogHandler *controllers.BlogHandler,
	aiHandler *controllers.AIHandler,
	oauthHandler *controllers.OAuthHandler,
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {
	router := gin.Default()

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// no auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/refresh", userHandler.RefreshToken)
			// Email and Password routes
			auth.POST("/send-verification", userHandler.SendVerificationEmail)
			auth.GET("/verify-email", userHandler.VerifyEmail)
			auth.POST("/forgot-password", userHandler.SendPasswordResetEmail)
			auth.POST("/reset-password", userHandler.ResetPassword)
			//Oauth routes
			auth.GET("/:provider/login", oauthHandler.OAuthLogin)
			auth.GET("/:provider/callback", oauthHandler.OAuthCallback)
		}
		// protected auth routes
		authProtected := v1.Group("/auth")
		authProtected.Use(authMiddleware.AuthRequired())
		{
			authProtected.POST("/logout", userHandler.Logout)
		}

		// user routes (authenticated)
		users := v1.Group("/users")
		users.Use(authMiddleware.AuthRequired())
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.POST("/profile/picture", userHandler.UploadProfilePicture)
		}
		// admin only routes
		admin := v1.Group("/admin")
		admin.Use(authMiddleware.AuthRequired(), authMiddleware.AdminRequired())
		{
			admin.PUT("/users/:id/promote", userHandler.PromoteUser)
			admin.PUT("/users/:id/demote", userHandler.DemoteUser)
		}
		// blog routes
		blogs := v1.Group("/blogs")
		{
			// public routes (no auth)
			blogs.GET("/", blogHandler.GetAllBlogs)
			blogs.GET("/:id", blogHandler.GetBlog)
			blogs.GET("/popular", blogHandler.GetPopularBlogs)

			//search and filter routes
			search := blogs.Group("/search")
			{
				search.GET("/title", blogHandler.SearchBlogsByTitle)
				search.GET("/author", blogHandler.SearchBlogsByAuthor)
			}

			filter := blogs.Group("/filter")
			{
				filter.GET("/tags", blogHandler.FilterBlogsByTags)
				filter.GET("/date", blogHandler.FilterBlogsByDate)
			}

			// protected routes (auth required)
			blogs.Use(authMiddleware.AuthRequired())
			blogs.POST("/", blogHandler.CreateBlog)
			blogs.PUT("/:id", blogHandler.UpdateBlog)
			blogs.DELETE("/:id", blogHandler.DeleteBlog)

			//comments

			blogs.POST("/:id/comments", blogHandler.AddComment)
			blogs.PUT("/:id/comments/:commentId", blogHandler.UpdateComment)
			blogs.DELETE("/:id/comments/:commentId", blogHandler.DeleteComment)

			//Reactions
			blogs.POST("/:id/like", blogHandler.LikeBlog)
			blogs.POST("/:id/dislike", blogHandler.DislikeBlog)
		}

		// AI routes (authenticated)
		ai := v1.Group("/ai")
		ai.Use(authMiddleware.AuthRequired())
		{
			ai.POST("/generate-blog", aiHandler.GenerateBlog)
			ai.POST("/enhance-blog", aiHandler.EnhanceBlog)
			ai.POST("/suggest-ideas", aiHandler.SuggestBlogIdeas)
		}
	}

	return router
}
