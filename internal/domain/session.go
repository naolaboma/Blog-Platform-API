package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// user session
type Session struct {
	ID                         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID                     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Username                   string             `bson:"username" json:"username"`
	Token                      string             `bson:"token" json:"token,omitempty"` // For JWT refresh token
	VerificationToken          string             `bson:"verification_token,omitempty" json:"verification_token,omitempty"`
	PasswordResetToken         string             `bson:"password_reset_token,omitempty" json:"password_reset_token,omitempty"`
	ResetCode                  int                `bson:"reset_code,omitempty" json:"reset_code,omitempty"`
	IsActive                   bool               `bson:"is_active" json:"is_active"`
	CreatedAt                  time.Time          `bson:"created_at" json:"created_at"`
	ExpiresAt                  time.Time          `bson:"expires_at" json:"expires_at"` // For JWT session
	VerificationTokenExpiresAt time.Time          `bson:"verification_expires_at,omitempty" json:"verification_expires_at,omitempty"`
	ResetTokenExpiresAt        time.Time          `bson:"reset_expires_at,omitempty" json:"reset_expires_at,omitempty"`
	LastActivity               time.Time          `bson:"last_activity" json:"last_activity"`
}

// interface for session data operations
type SessionRepository interface {
	Create(session *Session) error
	GetByID(id primitive.ObjectID) (*Session, error)
	GetByUserID(userID primitive.ObjectID) (*Session, error)
	GetByUsername(username string) (*Session, error)
	Update(session *Session) error
	Delete(id primitive.ObjectID) error
	DeleteByUserID(userID primitive.ObjectID) error
	DeleteExpired() error
	UpdateLastActivity(id primitive.ObjectID) error
	GetByVerificationToken(token string) (*Session, error)
	GetByResetToken(token string) (*Session, error)
}

// interface for session business logic
type SessionUseCase interface {
	CreateSession(userID primitive.ObjectID, username string, refreshToken string) (*Session, error)
	GetSessionByUserID(userID primitive.ObjectID) (*Session, error)
	DeleteSession(userID primitive.ObjectID) error
	CleanupExpiredSessions() error
	UpdateSessionActivity(userID primitive.ObjectID) error
}
