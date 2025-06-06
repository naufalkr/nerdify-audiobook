package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

const (
	// Verification token settings
	VerificationTokenLength = 32
	VerificationTokenExpiry = 24 * time.Hour
)

// GenerateVerificationToken creates a secure random token for email verification links
func GenerateVerificationToken() (string, error) {
	b := make([]byte, VerificationTokenLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// BuildVerificationLink creates the full verification URL with the token
func BuildVerificationLink(baseURL, token, email string) string {
	return fmt.Sprintf("%s/api/auth/verify-email-link?token=%s&email=%s", baseURL, token, email)
}

// BuildOTPAndLinkVerificationEmail creates an email with both OTP and verification link
func BuildOTPAndLinkVerificationEmail(fullName, otp, verificationLink string) string {
	// Read the email template file
	templateContent, err := ReadEmailTemplate("lecsens-email-template-dark.html")
	if err != nil {
		// Log the error
		fmt.Printf("Error reading email template: %v\n", err)

		// Fallback to simple template if the template file cannot be read
		return buildSimpleOTPEmail(fullName, otp, verificationLink)
	}

	// Replace placeholders with actual values
	templateContent = replacePlaceholder(templateContent, "{{.FullName}}", fullName)
	templateContent = replacePlaceholder(templateContent, "{{.OTP}}", otp)
	templateContent = replacePlaceholder(templateContent, "{{.VerifyLink}}", verificationLink)

	return templateContent
}

// buildSimpleOTPEmail creates a simple HTML email for OTP verification with a link
func buildSimpleOTPEmail(name, otp, verificationLink string) string {
	return `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Verifikasi Email LecSens</title>
		</head>
		<body>
			<h2>Halo, ` + name + `!</h2>
			<p>Terima kasih telah mendaftar. Berikut adalah kode verifikasi Anda:</p>
			<h1 style="font-size: 32px; letter-spacing: 5px; background-color: #f5f5f5; padding: 10px; text-align: center;">` + otp + `</h1>
			<p>Kode ini akan kadaluarsa dalam 15 menit.</p>
			<p>Atau, klik link berikut untuk memverifikasi email Anda secara langsung:</p>
			<p><a href="` + verificationLink + `" style="display: inline-block; padding: 10px 20px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 5px;">Verifikasi Email</a></p>
			<p>Jika Anda tidak mendaftar untuk akun ini, Anda dapat mengabaikan email ini.</p>
			<p>Terima kasih,</p>
			<p>Tim LecSens</p>
		</body>
		</html>
	`
}

// ReadEmailTemplate reads the content of an email template file
func ReadEmailTemplate(templateName string) (string, error) {
	// Get executable directory for more reliable path resolution
	execDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	// Build absolute path
	templatePath := filepath.Join(execDir, "helpers", "utils", "email_template", templateName)

	// Log the template path for debugging
	fmt.Printf("Reading email template from: %s\n", templatePath)

	// Read the template file
	return ReadFile(templatePath)
}

func SendMail(toEmail string, subject string, body string, urlPdf string) error {
	// Get email configuration from environment variables
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	senderName := os.Getenv("SMTP_SENDER_NAME")
	authEmail := os.Getenv("SMTP_AUTH_EMAIL")
	authPassword := os.Getenv("SMTP_AUTH_PASSWORD")

	// Convert port string to integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", senderName)
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	if urlPdf != "" {
		mailer.Attach(urlPdf)
	}

	dialer := gomail.NewDialer(
		host,
		port,
		authEmail,
		authPassword,
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}
