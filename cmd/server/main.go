package main

import (
	"fmt"
	"log"
	"net/http"

	"Blog-API/internal/delivery/controllers"
	"Blog-API/internal/delivery/router"
	"Blog-API/internal/infrastructure/ai"
	"Blog-API/internal/infrastructure/database"
	"Blog-API/internal/infrastructure/email"
	"Blog-API/internal/infrastructure/jwt"
	"Blog-API/internal/infrastructure/middleware"
	"Blog-API/internal/infrastructure/password"
	"Blog-API/internal/repository"
	"Blog-API/internal/usecase"
	"Blog-API/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	log.Printf("!!! DEBUG !!! Loaded Groq API Key: [%s]", cfg.AI.GroqAPIKey)

	gin.SetMode(cfg.Server.GinMode)

	mongoDB, err := database.NewMongoDB(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongoDB.Close()

	passwordService := password.NewPasswordService()
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)
	aiService := ai.NewAIService(cfg.AI.GroqAPIKey)
	baseURL := fmt.Sprintf("http://localhost:%s", cfg.Server.Port)
	emailService := email.NewEmailService(
		cfg.Email.Username,
		cfg.Email.Password,
		cfg.Email.Host,
		cfg.Email.Port,
		baseURL,
		cfg.Email.TemplatePath,
	)

	userRepo := repository.NewUserRepository(mongoDB)
	blogRepo := repository.NewBlogRepository(mongoDB)
	sessionRepo := repository.NewSessionRepository(mongoDB)

	userUseCase := usecase.NewUserUseCase(userRepo, passwordService, jwtService, sessionRepo, emailService)
	blogUseCase := usecase.NewBlogUseCase(blogRepo, userRepo)
	aiUseCase := usecase.NewAIUseCase(aiService)

	userHandler := controllers.NewUserHandler(userUseCase)
	blogHandler := controllers.NewBlogHandler(blogUseCase)
	aiHandler := controllers.NewAIHandler(aiUseCase)

	authMiddleware := middleware.NewAuthMiddleware(jwtService, sessionRepo)

	router := router.SetupRouter(userHandler, blogHandler, aiHandler, authMiddleware)

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Printf("MongoDB connected to: %s", cfg.MongoDB.URI)
	log.Printf("Gin mode: %s", cfg.Server.GinMode)

	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
