package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Blog-API/internal/delivery/controllers"
	"Blog-API/internal/delivery/router"
	"Blog-API/internal/infrastructure/ai"
	"Blog-API/internal/infrastructure/cache"
	"Blog-API/internal/infrastructure/database"
	"Blog-API/internal/infrastructure/email"
	"Blog-API/internal/infrastructure/filesystem"
	"Blog-API/internal/infrastructure/jwt"
	"Blog-API/internal/infrastructure/middleware"
	"Blog-API/internal/infrastructure/password"
	"Blog-API/internal/infrastructure/worker"
	"Blog-API/internal/repository"
	"Blog-API/internal/usecase"
	"Blog-API/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()
	log.Printf("!!! DEBUG !!! Loaded Groq API Key: [%s]", cfg.AI.GroqAPIKey)

	gin.SetMode(cfg.Server.GinMode)

	mongoDB, err := database.NewMongoDB(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	//defer mongoDB.Close()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	//---worker Pool---
	workerPool := worker.NewPool(4, 100)
	workerPool.Start()
	//---Services---
	passwordService := password.NewPasswordService()
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)
	aiService := ai.NewAIService(cfg.AI.GroqAPIKey)
	baseURL := fmt.Sprintf("http://localhost:%s", cfg.Server.Port)
	fileService := filesystem.NewFileService(cfg.Upload.Path)
	emailService := email.NewEmailService(
		cfg.Email.Username,
		cfg.Email.Password,
		cfg.Email.Host,
		cfg.Email.Port,
		baseURL,
		cfg.Email.TemplatePath,
	)
	cacheService := cache.NewRedisCache(redisClient)
	//---repositories---
	userRepo := repository.NewUserRepository(mongoDB)
	blogRepo := repository.NewBlogRepository(mongoDB)
	sessionRepo := repository.NewSessionRepository(mongoDB)
	//---use cases---
	userUseCase := usecase.NewUserUseCase(userRepo, passwordService, jwtService, sessionRepo, emailService, fileService, workerPool)
	blogUseCase := usecase.NewBlogUseCase(blogRepo, userRepo, cacheService)
	aiUseCase := usecase.NewAIUseCase(aiService)
	//---handlers---
	userHandler := controllers.NewUserHandler(userUseCase)
	blogHandler := controllers.NewBlogHandler(blogUseCase)
	aiHandler := controllers.NewAIHandler(aiUseCase)

	authMiddleware := middleware.NewAuthMiddleware(jwtService, sessionRepo)
	router := router.SetupRouter(userHandler, blogHandler, aiHandler, authMiddleware)

	//Graceful server shutdown logic S

	httpServer := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	// wait for an interrupt signal (like ctrl + c)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// shutdown the http server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	workerPool.Shutdown()
	if err := redisClient.Close(); err != nil {
		log.Printf("Failed to close Redis client: %v", err)
	}
	//close the mongodb connect
	mongoDB.Close()
	log.Println("Server exiting.")
}
