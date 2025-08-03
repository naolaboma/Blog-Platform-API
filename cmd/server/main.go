package main

import (
	"log"
	"net/http"

	"blog-api-mongodb/internal/delivery/controllers"
	"blog-api-mongodb/internal/delivery/router"
	"blog-api-mongodb/internal/infrastructure/database"
	"blog-api-mongodb/internal/infrastructure/jwt"
	"blog-api-mongodb/internal/infrastructure/middleware"
	"blog-api-mongodb/internal/infrastructure/password"
	"blog-api-mongodb/internal/repository"
	"blog-api-mongodb/internal/usecase"
	"blog-api-mongodb/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	gin.SetMode(cfg.Server.GinMode)

	mongoDB, err := database.NewMongoDB(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongoDB.Close()

	passwordService := password.NewPasswordService()
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)

	userRepo := repository.NewUserRepository(mongoDB)
	sessionRepo := repository.NewSessionRepository(mongoDB)

	userUseCase := usecase.NewUserUseCase(userRepo, passwordService, jwtService, sessionRepo)

	userHandler := controllers.NewUserHandler(userUseCase)

	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	router := router.SetupRouter(userHandler, authMiddleware)

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Printf("MongoDB connected to: %s", cfg.MongoDB.URI)
	log.Printf("Gin mode: %s", cfg.Server.GinMode)

	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 