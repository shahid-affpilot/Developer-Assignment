package email

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func LoadEmailConfig() *EmailConfig {
	return &EmailConfig{
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     os.Getenv("EMAIL_PORT"),
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		From:     os.Getenv("EMAIL_FROM"),
	}
}

func SendVerificationEmail(to, username, token string) error {
	config := LoadEmailConfig()
	verificationURL := fmt.Sprintf("%s?token=%s",
		os.Getenv("EMAIL_VERIFICATION_URL"), token)

	subject := "Verify Your Email Address"
	body := fmt.Sprintf(`
Hello %s,

Thank you for registering. Please verify your email address by clicking the link below:

%s

This link will expire in 5 minutes.

Best regards,
AffPilot Team
`, username, verificationURL)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", config.From, to, subject, body)

	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	return smtp.SendMail(addr, auth, config.From, []string{to}, []byte(msg))
}
