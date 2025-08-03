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
	// Validate password
	if err := u.passwordService.ValidatePassword(password); err != nil {
		return nil, err
	}

	// Check if user already exists by email
	existingUser, _ := u.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Check if user already exists by username
	existingUser, _ = u.userRepo.GetByUsername(username)
	if existingUser != nil {
		return nil, errors.New("user with this username already exists")
	}

	// Hash password
	hashedPassword, err := u.passwordService.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &domain.User{
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		Role:      "user", // Default role
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) Login(email, password string) (*domain.LoginResponse, error) {
	// Get user by email
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if !u.passwordService.CheckPassword(password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// TODO:Implement JWT login business logic
	// Requirements:
	// - Get user by email using userRepo.GetByEmail(email) - already done above
	// - Check password using passwordService.CheckPassword(password, user.Password) - already done above
	// - Generate access token using jwtService.GenerateAccessToken(user.ID, user.Email, user.Role)
	// - Generate refresh token using jwtService.GenerateRefreshToken(user.ID, user.Email, user.Role)
	// - Create session with refresh token using sessionRepo.Create()
	// - Return LoginResponse with user, access_token, and refresh_token
	return nil, nil
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
	// RefreshToken refreshes an access token using a refresh token
	// TODO: Implement token refresh business logic
	// Requirements:
	// - Validate refresh token using jwtService.ValidateToken(refreshToken)
	// - Get session by userID using sessionRepo.GetByUserID(claims.UserID)
	// - Check if session is active and not expired
	// - Get user by ID using userRepo.GetByID(claims.UserID)
	// - Generate new access token using jwtService.GenerateAccessToken()
	// - Update session activity using sessionRepo.UpdateLastActivity()
	// - Return LoginResponse with user, new access_token, and same refresh_token
	return nil, nil
}

func (u *UserUseCase) Logout(userID primitive.ObjectID) error {
	// TODO: Implement logout business logic
	// Requirements:
	// - Delete session from database using sessionRepo.DeleteByUserID(userID)
	// - Return error if any
	return nil
}
