package usecase

import (
	"Blog-API/internal/domain"
	"context"
)

type EmailJob struct {
	EmailService           domain.EmailService
	Type                   string
	Email, Username, Token string
}

func (j *EmailJob) Run(ctx context.Context) error {
	switch j.Type {
	case "verification":
		return j.EmailService.SendVerificationEmail(j.Email, j.Username, j.Token)
	case "password_reset":
		return j.EmailService.SendPasswordResetEmail(j.Email, j.Username, j.Token)
	}
	return nil

}
