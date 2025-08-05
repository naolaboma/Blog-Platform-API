package router

import (
	"Blog-API/internal/delivery/controllers"
	"Blog-API/internal/infrastructure/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *controllers.UserHandler, blogHandler *controllers.BlogHandler, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
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
		}

		// blog routes (authenticated)
		blogs := v1.Group("/blogs")
		blogs.Use(authMiddleware.AuthRequired())
		{
			blogs.POST("/", blogHandler.CreateBlog)
			blogs.GET("/:id", blogHandler.GetBlog)
			blogs.PUT("/:id", blogHandler.UpdateBlog)
			blogs.DELETE("/:id", blogHandler.DeleteBlog)
		}

		// blog search routes (public)
		blogsPublic := v1.Group("/blogs")
		{
			blogsPublic.GET("/search/title", blogHandler.SearchBlogsByTitle)
			blogsPublic.GET("/search/author", blogHandler.SearchBlogsByAuthor)
		}
	}

	return router
}
