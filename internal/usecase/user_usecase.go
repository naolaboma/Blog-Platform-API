package usecase

import (
	"Blog-API/internal/domain"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUseCase struct {
	userRepo        domain.UserRepository
	passwordService domain.PasswordService
	jwtService      domain.JWTService
	sessionRepo     domain.SessionRepository
	emailService    domain.EmailService
	fileService     domain.FileService
	workerPool      domain.WorkerPool
}

func NewUserUseCase(
	userRepo domain.UserRepository,
	passwordService domain.PasswordService,
	jwtService domain.JWTService,
	sessionRepo domain.SessionRepository,
	emailService domain.EmailService,
	fileService domain.FileService,
	workerPool domain.WorkerPool,
) domain.UserUseCase {
	return &UserUseCase{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
		sessionRepo:     sessionRepo,
		emailService:    emailService,
		fileService:     fileService,
		workerPool:      workerPool,
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

	// Generate access token
	accessToken, err := u.jwtService.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := u.jwtService.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	// Create session with refresh token
	session := &domain.Session{
		UserID:       user.ID,
		Username:     user.Username,
		Token:        refreshToken,
		IsActive:     true,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(time.Hour * 24 * 7), // exp in 7 days
		LastActivity: time.Now(),
	}
	if err := u.sessionRepo.Create(session); err != nil {
		return nil, err
	}

	// Return login response
	return &domain.LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *UserUseCase) GetByID(id primitive.ObjectID) (*domain.User, error) {
	return u.userRepo.GetByID(id)
}

func (u *UserUseCase) UpdateProfile(id primitive.ObjectID, req *domain.UpdateProfileRequest) (*domain.User, error) {
	// Check if user exists
	currentUser, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Prepare updates
	updates := make(map[string]interface{})

	// Handle username update with uniqueness check
	if req.Username != nil && *req.Username != currentUser.Username {
		// Check if username already exists (excluding current user)
		existingUser, _ := u.userRepo.GetByUsername(*req.Username)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("username already exists")
		}
		updates["username"] = *req.Username
	}

	// Handle email update with uniqueness check
	if req.Email != nil && *req.Email != currentUser.Email {
		// Check if email already exists (excluding current user)
		if existingUser, _ := u.userRepo.GetByEmail(*req.Email); existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already exists")
		}
		updates["email"] = *req.Email
	}

	// Handle bio update
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}

	// Add updated_at timestamp
	updates["updated_at"] = time.Now()

	// Update user if there are any changes
	if len(updates) > 1 { // More than just updated_at
		if err := u.userRepo.UpdateProfile(id, updates); err != nil {
			return nil, err
		}
	}

	// Return updated user
	return u.userRepo.GetByID(id)
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
	claims, err := u.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get the corresponding session from the database
	session, err := u.sessionRepo.GetByUserID(claims.UserID)
	if err != nil {
		return nil, errors.New("session not found")
	}

	// Check if the session is still active and has not expired
	if !session.IsActive || time.Now().After(session.ExpiresAt) {
		return nil, errors.New("session is expired or inactive")
	}

	// Get the full user details
	user, err := u.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new access token
	newAccessToken, err := u.jwtService.GenerateAccessToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	// Update last active time on the session
	err = u.sessionRepo.UpdateLastActivity(session.ID)
	if err != nil {
		return nil, err
	}

	// Return the response with the new access token and the original refresh token
	return &domain.LoginResponse{
		User:         user,
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *UserUseCase) Logout(userID primitive.ObjectID) error {
	return u.sessionRepo.DeleteByUserID(userID)
}

func (u *UserUseCase) VerifyEmail(token string) error {
	// find the session associated with this token
	session, err := u.sessionRepo.GetByVerificationToken(token)
	if err != nil {
		return errors.New("invalid or expired verification token")
	}
	// check if the token has expired
	if time.Now().After(session.VerificationTokenExpiresAt) {
		return errors.New("invalid or expired verification token")
	}

	// update/mark user's email as verified in the database
	if err := u.userRepo.UpdateEmailVerificationStatus(session.UserID, true); err != nil {
		return err
	}
	// clean verif token from session so it cant be reused
	session.VerificationToken = ""
	return u.sessionRepo.Update(session)
}

func (u *UserUseCase) SendVerificationEmail(email string) error {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return errors.New("if a user with this email exists, a verification email has been sent")
	}

	// check if already verified
	if user.EmailVerified {
		return errors.New("email is already verified")
	}
	// generate secure, short lived verification token
	verificationToken := u.passwordService.GenerateSecureToken(32)

	//store the token and its expiry in a session document
	//(we reuse the session logic for simplicity)

	session := &domain.Session{
		UserID:                     user.ID,
		Username:                   user.Username,
		VerificationToken:          verificationToken,
		VerificationTokenExpiresAt: time.Now().Add(24 * time.Hour),
		IsActive:                   false, // this is not a login session
	}
	// we create or update session entry for this user
	existingSession, _ := u.sessionRepo.GetByUserID(user.ID)
	if existingSession != nil {
		existingSession.VerificationToken = session.VerificationToken
		existingSession.VerificationTokenExpiresAt = session.VerificationTokenExpiresAt
		if err := u.sessionRepo.Update(existingSession); err != nil {
			return err
		}
	} else {
		if err := u.sessionRepo.Create(session); err != nil {
			return err
		}
	}

	// send the email in Background
	u.workerPool.Submit(&EmailJob{
		EmailService: u.emailService,
		Type:         "verification",
		Email:        user.Email,
		Username:     user.Username,
		Token:        verificationToken,
	})
	return nil
}

func (u *UserUseCase) SendPasswordResetEmail(email string) error {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return errors.New("if a user with this email exists, a password reset email has been sent")
	}

	resetToken := u.passwordService.GenerateSecureToken(32)

	session, err := u.sessionRepo.GetByUserID(user.ID)
	if err != nil || session == nil {
		session = &domain.Session{UserID: user.ID, Username: user.Username}
	}

	session.PasswordResetToken = resetToken
	session.ResetTokenExpiresAt = time.Now().Add(1 * time.Hour)
	if err != nil {
		if err := u.sessionRepo.Create(session); err != nil {
			return err
		}
	} else {
		if err := u.sessionRepo.Update(session); err != nil {
			return err
		}
	}

	//go u.emailService.SendPasswordResetEmail(user.Email, user.Username, resetToken)
	u.workerPool.Submit(&EmailJob{
		EmailService: u.emailService,
		Type:         "password_reset",
		Email:        user.Email,
		Username:     user.Username,
		Token:        resetToken,
	})
	return nil
}

func (u *UserUseCase) ResetPassword(token, newPassword string) error {
	// validate the new pasword's strength
	if err := u.passwordService.ValidatePassword(newPassword); err != nil {
		return err
	}

	// find the session associated with the reset token
	session, err := u.sessionRepo.GetByResetToken(token)
	if err != nil {
		return errors.New("invalid or expired passoword reset token")
	}
	if time.Now().After(session.ResetTokenExpiresAt) {
		return errors.New("invalid or expired password reset token")
	}
	// hash the new password
	hashedPassword, err := u.passwordService.HashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := u.userRepo.UpdatePassword(session.UserID, hashedPassword); err != nil {
		return err
	}
	session.PasswordResetToken = ""
	return u.sessionRepo.Update(session)
}

//Updated the updaterole and added the upload profile picture functions

func (u *UserUseCase) UpdateRole(adminUserID, targetUserID primitive.ObjectID, role string) error {
	adminUser, err := u.userRepo.GetByID(adminUserID)
	if err != nil {
		return errors.New("admin user not found")
	}
	if adminUser.Role != domain.RoleAdmin {
		return errors.New("target user not found")
	}
	if adminUserID == targetUserID && role == domain.RoleUser {
		return errors.New("admins cannot demote themselves")
	}
	return u.userRepo.UpdateRole(targetUserID, role)
}
func (u *UserUseCase) UploadProfilePicture(userID primitive.ObjectID, file multipart.File, handler *multipart.FileHeader) (*domain.User, error) {
	// 1. Save the file using the file service interface.
	photo, err := u.fileService.SaveProfilePicture(userID, file, handler)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// 2. Update the user's document in the database.
	if err := u.userRepo.UpdateProfilePicture(userID, photo); err != nil {
		return nil, fmt.Errorf("failed to update user profile in database: %w", err)
	}

	// 3. Return the fully updated user object.
	return u.userRepo.GetByID(userID)
}
