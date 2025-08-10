package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server  ServerConfig
	MongoDB MongoDBConfig
	JWT     JWTConfig
	Email   EmailConfig
	Upload  UploadConfig
	AI      AIConfig
	Redis   RedisConfig
	OAuth   OAuthConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type EmailConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	From         string
	TemplatePath string
}

type UploadConfig struct {
	Path        string
	MaxFileSize int64
}
type AIConfig struct {
	GroqAPIKey string
}
type RedisConfig struct {
	Addr     string `mapstructure:"REDIS_ADDR"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}
type OAuthProvider struct {
	ClientID     string   `mapstructure:"CLIENT_ID"`
	ClientSecret string   `mapstructure:"CLIENT_SECRET"`
	RedirectURL  string   `mapstructure:"REDIRECT_URL"`
	Scopes       []string `mapstructure:"SCOPES"`
}
type OAuthConfig struct {
	Google      OAuthProvider
	GitHub      OAuthProvider
	StateSecret string `mapstructure:"OAUTH_STATE_SECRET"`
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Server: ServerConfig{
			Port:    getEnv("SERVER_PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017/blog_db"),
			Database: getEnv("MONGODB_DATABASE", "blog_db"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "your-super-secret-jwt-key-here"),
			AccessExpiry:  getDurationEnv("JWT_ACCESS_EXPIRY", 15*time.Minute),
			RefreshExpiry: getDurationEnv("JWT_REFRESH_EXPIRY", 168*time.Hour), // 7 days
		},
		Email: EmailConfig{
			Host:         getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:         getIntEnv("SMTP_PORT", 587),
			Username:     getEnv("SMTP_EMAIL", ""),
			Password:     getEnv("SMTP_PASSWORD", ""),
			From:         getEnv("EMAIL_FROM", "noreply@blogplatform.com"),
			TemplatePath: getEnv("EMAIL_TEMPLATE_PATH", "./internal/infrastructure/email/templates"),
		},
		Upload: UploadConfig{
			Path:        getEnv("UPLOAD_PATH", "./uploads"),
			MaxFileSize: getInt64Env("MAX_FILE_SIZE", 5*1024*1024), // 5MB
		},
		AI: AIConfig{
			GroqAPIKey: getEnv("GROQ_API_KEY", ""),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		OAuth: OAuthConfig{
			StateSecret: getEnv("OAUTH_STATE_SECRET", "a-very-secret-string-for-oauth-state-change-me"),
			Google: OAuthProvider{
				ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback"),
				Scopes:       getScopes("GOOGLE_SCOPES", "https://www.googleapis.com/auth/userinfo.email,https://www.googleapis.com/auth/userinfo.profile"),
			},
			GitHub: OAuthProvider{
				ClientID:     getEnv("GITHUB_CLIENT_ID", ""),
				ClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GITHUB_REDIRECT_URL", "http://localhost:8080/api/v1/auth/github/callback"),
				Scopes:       getScopes("GITHUB_SCOPES", "read:user,user:email"),
			},
		},
	}
}

// Helper functions to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
func getScopes(key, defaultValue string) []string {
	value := getEnv(key, defaultValue)
	return strings.Split(value, ",")
}
