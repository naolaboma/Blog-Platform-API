package email

import (
	"Blog-API/internal/domain"
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"path/filepath"
)

type EmailService struct {
	auth        smtp.Auth
	from        string
	host        string
	port        int
	baseURL     string
	templateDir string
}

type EmailData struct {
	Username string
	Token    string
	Link     string
	Subject  string
	To       string
}

// type EmailTemplate struct {
// 	Subject string
// 	Body    string
// }

//	func NewEmailService() *EmailService {
//		auth := smtp.PlainAuth(
//			"",
//			os.Getenv("SMTP_EMAIL"),
//			os.Getenv("SMTP_PASSWORD"),
//			"smtp.gmail.com",
//		)
//		return &EmailService{auth: auth}
//	}
func NewEmailService(username, password, host string, port int, baseURL, templatePath string) domain.EmailService {
	auth := smtp.PlainAuth("", username, password, host)
	return &EmailService{
		auth:        auth,
		from:        username,
		host:        host,
		port:        port,
		baseURL:     baseURL,
		templateDir: templatePath,
	}
}

func (e *EmailService) SendVerificationEmail(to, username, token string) error {
	data := EmailData{
		Username: username,
		Token:    token,
		Link:     fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", e.baseURL, token),
		Subject:  "Verify Your Email Address",
		To:       to,
	}

	return e.sendEmail("verification.html", data)
}

func (e *EmailService) SendPasswordResetEmail(to, username, token string) error {
	data := EmailData{
		Username: username,
		Token:    token,
		Link:     fmt.Sprintf("%s/reset-password?token=%s", e.baseURL, token),
		Subject:  "Reset Your Password",
		To:       to,
	}

	return e.sendEmail("password_reset.html", data)
}

func (e *EmailService) sendEmail(templateName string, data EmailData) error {
	// Load and parse base + content templates
	tmplt, err := template.ParseFiles(
		filepath.Join(e.templateDir, "base.html"),
		filepath.Join(e.templateDir, templateName),
	)
	if err != nil {
		return fmt.Errorf("error loading templates: %v", err)
	}

	// Execute base template which includes the named "content" block
	var bodyWritten bytes.Buffer
	if err := tmplt.ExecuteTemplate(&bodyWritten, "base.html", data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	// Create email
	from := e.from
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n%s", from, data.To, data.Subject, bodyWritten.String())

	// Send email using configured host:port
	addr := fmt.Sprintf("%s:%d", e.host, e.port)
	if err := smtp.SendMail(addr, e.auth, from, []string{data.To}, []byte(msg)); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	return nil
}

func (e *EmailService) SendWelcomeEmail(email, username string) error {
	return nil
}
