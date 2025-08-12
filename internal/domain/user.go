package domain

import (
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username       string             `bson:"username" json:"username" validate:"required,min=3,max=50"`
	Email          string             `bson:"email" json:"email" validate:"required,email"`
	Password       string             `bson:"password" json:"-" validate:"required,min=6"` // "-" means don't include in JSON
	Role           string             `bson:"role" json:"role"`
	EmailVerified  bool               `bson:"email_verified" json:"email_verified"`
	ProfilePicture *Photo             `bson:"profile_picture,omitempty" json:"profile_picture,omitempty"`
	Bio            string             `bson:"bio,omitempty" json:"bio,omitempty"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
	OAuthProvider  string             `bson:"oauth_provider,omitempty" json:"oauth_provider,omitempty"`
	OAuthID        string             `bson:"oauth_id,omitempty" json:"oauth_id,omitempty"`
}

// Constants for user roles to avoid magic strings.
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type Photo struct {
	Filename   string    `bson:"filename" json:"filename"`
	FilePath   string    `bson:"file_path" json:"file_path"`
	PublicID   string    `bson:"public_id" json:"public_id"`
	UploadedAt time.Time `bson:"uploaded_at" json:"uploaded_at"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id primitive.ObjectID) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(user *User) error
	Delete(id primitive.ObjectID) error
	UpdateProfile(id primitive.ObjectID, updates map[string]interface{}) error
	UpdatePassword(id primitive.ObjectID, password string) error
	UpdateRole(id primitive.ObjectID, role string) error
	UpdateProfilePicture(id primitive.ObjectID, photo *Photo) error
	VerifyEmail(id primitive.ObjectID) error
	UpdateEmailVerificationStatus(id primitive.ObjectID, verified bool) error

	GetByOAuth(provider, oauthID string) (*User, error)
}

type UserUseCase interface {
	Register(username, email, password string) (*User, error)
	Login(email, password string) (*LoginResponse, error)
	GetByID(id primitive.ObjectID) (*User, error)
	ValidatePassword(password string) error
	HashPassword(password string) (string, error)
	CheckPassword(password, hash string) bool
	RefreshToken(refreshToken string) (*LoginResponse, error)
	Logout(userID primitive.ObjectID) error
	VerifyEmail(token string) error
	SendVerificationEmail(email string) error
	SendPasswordResetEmail(email string) error
	ResetPassword(token, newPassword string) error

	UpdateProfile(id primitive.ObjectID, req *UpdateProfileRequest) (*User, error)
	UpdateRole(adminUserID, targetUserID primitive.ObjectID, role string) error
	UploadProfilePicture(userID primitive.ObjectID, file multipart.File, handler *multipart.FileHeader) (*User, error)

	OAuthLogin(provider, state, code, storedState string) (*LoginResponse, error)
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UpdateProfileRequest struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Bio      *string `json:"bio,omitempty" validate:"omitempty,max=500"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type EmailVerificationRequest struct {
	Token string `json:"token" validate:"required"`
}

type EmailVerificationResponse struct {
	Message string `json:"message"`
	User    *User  `json:"user"`
}

type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordResetResponse struct {
	Message string `json:"message"`
}

type NewPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type NewPasswordResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=user admin"`
}
