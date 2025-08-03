package usecase

import (
	"Blog-API/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionUseCase struct {
	sessionRepo domain.SessionRepository
}

func NewSessionUseCase(sessionRepo domain.SessionRepository) domain.SessionUseCase {
	return &SessionUseCase{
		sessionRepo: sessionRepo,
	}
}

func (s *SessionUseCase) CreateSession(userID primitive.ObjectID, username string, refreshToken string) (*domain.Session, error) {

	// TODO: Implement session creation business logic
	// Requirements:
	// - Create a new Session object with provided userID and username
	// - Set IsActive to true
	// - Set ExpiresAt to 7 days from now (time.Now().Add(7 * 24 * time.Hour))
	// - Call sessionRepo.Create() to save to database
	// - Return the created session or error

	session := &domain.Session{
		UserID: userID,
		Username: username,
		Token: refreshToken, // store the refreshToken
		IsActive: true,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // exp for 7 days
		LastActivity: time.Now(),
	}
	if err := s.sessionRepo.Create(session); err != nil{
		return nil, err
	}
	return session, nil
}

func (s *SessionUseCase) GetSessionByUserID(userID primitive.ObjectID) (*domain.Session, error) {
	// TODO: Implement session retrieval business logic
	// Requirements:
	// - Call sessionRepo.GetByUserID(userID)
	// - Return the session or error
	return s.sessionRepo.
}

func (s *SessionUseCase) DeleteSession(userID primitive.ObjectID) error {
	// TODO:Implement session deletion business logic
	// Requirements:
	// - Call sessionRepo.DeleteByUserID(userID)
	// - Return error if any
	return nil
}

func (s *SessionUseCase) CleanupExpiredSessions() error {
	// TODO: Implement expired session cleanup business logic
	// Requirements:
	// - Call sessionRepo.DeleteExpired()
	// - Return error if any
	return nil
}

func (s *SessionUseCase) UpdateSessionActivity(userID primitive.ObjectID) error {
	// TODO: Implement session activity update business logic
	// Requirements:
	// - First get the session by userID using GetSessionByUserID
	// - If session exists, call sessionRepo.UpdateLastActivity(session.ID)
	// - Return error if any
	return nil
}
