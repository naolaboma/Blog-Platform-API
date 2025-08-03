package usecase

import (
	"errors"
	"time"

	"Blog-API/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUseCase struct {
	userRepo        domain.UserRepository
	passwordService domain.PasswordService
	jwtService      domain.JWTService
	sessionRepo     domain.SessionRepository
}

func NewUserUseCase(userRepo domain.UserRepository, passwordService domain.PasswordService, jwtService domain.JWTService, sessionRepo domain.SessionRepository) domain.UserUseCase {
	return &UserUseCase{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
		sessionRepo:     sessionRepo,
	}
}

func (u *UserUseCase) Register(username, email, password string) (*domain.User, error) {
	if err := u.passwordService.ValidatePassword(password); err != nil {
		return nil, err
	}

	existingUser, _ := u.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	existingUser, _ = u.userRepo.GetByUsername(username)
	if existingUser != nil {
		return nil, errors.New("user with this username already exists")
	}

	hashedPassword, err := u.passwordService.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		Role:      "user", // Default role
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) Login(email, password string) (*domain.LoginResponse, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !u.passwordService.CheckPassword(password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	accessToken, err := u.jwtService.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	
	refreshToken, err := u.jwtService.GenerateRefreshToken(user.Id, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	session := &domain.Session{
		UserID: user.ID,
		Username: user.Username,
		Token: refreshToken,
		IsActive: true,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7), // exp in 7 days
		LastActivity: time.Now(),
	}
	if err := u.sessionRepo.Create(session); err != nil{
		return nil, err
	}

	return &domain.LoginResponse{
		User: user,
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *UserUseCase) GetByID(id primitive.ObjectID) (*domain.User, error) {
	return u.userRepo.GetByID(id)
}

func (u *UserUseCase) UpdateProfile(id primitive.ObjectID, bio, profilePic, contactInfo *string) (*domain.User, error) {
	// Check if user exists
	_, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Prepare updates
	updates := make(map[string]interface{})
	if bio != nil {
		updates["bio"] = *bio
	}
	if profilePic != nil {
		updates["profile_picture"] = *profilePic
	}
	updates["updated_at"] = time.Now()

	// Update user
	if err := u.userRepo.UpdateProfile(id, updates); err != nil {
		return nil, err
	}

	// Return updated user
	return u.userRepo.GetByID(id)
}

func (u *UserUseCase) UpdateRole(id primitive.ObjectID, role string) error {
	if role != "user" && role != "admin" {
		return errors.New("invalid role")
	}
	return u.userRepo.UpdateRole(id, role)
}

func (u *UserUseCase) ValidatePassword(password string) error {
	return u.passwordService.ValidatePassword(password)
}

func (u *UserUseCase) HashPassword(password string) (string, error) {
	return u.passwordService.HashPassword(password)
}

func (u *UserUseCase) CheckPassword(password, hash string) bool {
	return u.passwordService.CheckPassword(password, hash)
}

func (u *UserUseCase) RefreshToken(refreshToken string) (*domain.LoginResponse, error) {
	claims, err := u.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	session, err := u.sessionRepo.GetByUserID(claims.UserID)
	if err != nil {
		return nil, errors.New("session not found")
	}

	if !session.IsActive || time.Now().After(session.ExpiresAt) {
		return nil, errors.New("session expired")
	}

	// Get user
	user, err := u.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	
	newAccessToken, err := u.jwtService.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	err = u.sessionRepo.UpdateLastActivity(session.ID)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		User:         user,
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *UserUseCase) Logout(userID primitive.ObjectID) error {
	return u.sessionRepo.DeleteByUserID(userID)
} 
