package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/shahid-affpilot/affpilot-auth-service/internal/config"
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

func SendVerificationEmail(to, subject, body string) error {
	cnf := config.GetConfig()
	emailConfig := cnf.Email

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", emailConfig.From, to, subject, body)

	auth := smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.Host)
	addr := fmt.Sprintf("%s:%s", emailConfig.Host, strconv.Itoa(emailConfig.Port))

	return smtp.SendMail(addr, auth, emailConfig.From, []string{to}, []byte(msg))
}

func GenerateVerificationURL(userID uuid.UUID) (string, error) {
	cnf := config.GetConfig()
	ttl := time.Duration(cnf.Email.VerificationTTL) * time.Minute
	now := time.Now()

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"created": now.Unix(),
		"expiry":  now.Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cnf.JWT.Secret))
	if err != nil {
		log.Printf("Failed to sign verification token: %v", err)
		return "", err
	}

	verificationURL := fmt.Sprintf("%s?token=%s", cnf.Email.VerificationURL, signedToken)
	return verificationURL, nil
}
