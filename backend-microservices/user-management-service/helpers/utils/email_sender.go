package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

// EmailSender defines the interface for sending emails
type EmailSender interface {
	Send(to string, subject string, content string) error
}

// SMTPEmailSender implements the EmailSender interface with SMTP
type SMTPEmailSender struct {
	config struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
}

// NewSMTPEmailSender creates a new SMTPEmailSender with explicit configuration
func NewSMTPEmailSender(host string, port int, username, password, senderName string) EmailSender {
	return &SMTPEmailSender{
		config: struct {
			Host     string
			Port     int
			Username string
			Password string
			Sender   string
		}{
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
			Sender:   fmt.Sprintf("%s <%s>", senderName, username),
		},
	}
}

// Send sends an email using SMTP
func (s *SMTPEmailSender) Send(to string, subject string, content string) error {
	// Create a new message
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.Sender)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	// Add additional headers for better deliverability
	m.SetHeader("X-Priority", "1")
	m.SetHeader("X-MSMail-Priority", "High")
	m.SetHeader("X-Mailer", "LecSens Mailer")
	m.SetHeader("List-Unsubscribe", "<mailto:"+s.config.Username+"?subject=unsubscribe>")
	m.SetHeader("Precedence", "bulk")

	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)
	return d.DialAndSend(m)
}

// BuildOTPVerificationEmail creates the HTML content for OTP verification email
func BuildOTPVerificationEmail(name, otp string) string {
	templatePath := filepath.Join("helpers", "utils", "email_template", "base-template.html")

	// Get current working directory
	execDir, err := os.Getwd()
	if err != nil {
		return buildBasicOTPEmail(name, otp)
	}

	// Build absolute path
	templatePath = filepath.Join(execDir, templatePath)

	// Log the template path for debugging
	fmt.Printf("Reading email template from: %s\n", templatePath)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		// Fallback to simple template if file cannot be parsed
		fmt.Printf("Error parsing template: %v\n", err)
		return buildBasicOTPEmail(name, otp)
	}

	// Get logo URL from environment variable or use default
	logoURL := os.Getenv("EMAIL_LOGO_URL")
	if logoURL == "" {
		// Always use the app URL for logo reference
		logoURL = os.Getenv("APP_URL") + "/static/logo.png"
	}

	data := struct {
		Name      string
		OTP       string
		VerifyURL string
		LogoURL   string
	}{
		Name:      name,
		OTP:       otp,
		VerifyURL: os.Getenv("APP_URL") + "/verify-email",
		LogoURL:   logoURL,
	}

	var content bytes.Buffer
	if err := tmpl.Execute(&content, data); err != nil {
		return buildBasicOTPEmail(name, otp)
	}

	return content.String()
}

// BuildPasswordResetEmail creates the HTML content for password reset email
func BuildPasswordResetEmail(name, token string) string {
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "https://lecsens-iot.erplabiim.com" // Default to production URL
	}
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", appURL, token)
	templatePath := filepath.Join("utils", "email_template", "lecsens-forget-password-template-dark.html")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		// Fallback to simple template if file cannot be parsed
		return buildSimpleResetEmail(name, resetURL)
	}

	data := struct {
		Name       string
		VerifyLink string
	}{
		Name:       name,
		VerifyLink: resetURL,
	}

	var content bytes.Buffer
	if err := tmpl.Execute(&content, data); err != nil {
		return buildSimpleResetEmail(name, resetURL)
	}

	return content.String()
}

// buildBasicOTPEmail creates a simple HTML email for OTP verification
func buildBasicOTPEmail(name, otp string) string {
	return `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Kode Verifikasi Email</title>
		</head>
		<body>
			<h2>Halo, ` + name + `!</h2>
			<p>Terima kasih telah mendaftar. Berikut adalah kode verifikasi Anda:</p>
			<h1 style="font-size: 32px; letter-spacing: 5px; background-color: #f5f5f5; padding: 10px; text-align: center;">` + otp + `</h1>
			<p>Kode ini akan kadaluarsa dalam 15 menit.</p>
			<p>Jika Anda tidak mendaftar untuk akun ini, Anda dapat mengabaikan email ini.</p>
			<p>Terima kasih,</p>
			<p>Tim LecSens</p>
		</body>
		</html>
	`
}

// buildSimpleResetEmail creates a simple HTML email for password reset
func buildSimpleResetEmail(name, resetURL string) string {
	return `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Reset Password</title>
		</head>
		<body>
			<h2>Halo, ` + name + `!</h2>
			<p>Kami menerima permintaan untuk mengatur ulang password Anda. Silakan klik tautan di bawah ini untuk mengatur password baru:</p>
			<p><a href="` + resetURL + `">Reset Password</a></p>
			<p>Tautan ini akan kadaluarsa dalam 1 jam.</p>
			<p>Jika Anda tidak meminta pengaturan ulang password, silakan abaikan email ini.</p>
			<p>Terima kasih,</p>
			<p>Tim LecSens</p>
		</body>
		</html>
	`
}

// BuildEmailUpdateOTPVerificationEmail creates the HTML content for email update verification
func BuildEmailUpdateOTPVerificationEmail(name, oldEmail, newEmail, otp string) string {
	templatePath := filepath.Join("utils", "email_template", "base-template-dark.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		// Fallback to simple template if file cannot be parsed
		return buildSimpleEmailUpdateOTPEmail(name, oldEmail, newEmail, otp)
	}

	// Get logo URL from environment variable or use default
	logoURL := os.Getenv("EMAIL_LOGO_URL")
	if logoURL == "" {
		// Always use the app URL for logo reference
		logoURL = os.Getenv("APP_URL") + "/static/logo.png"
	}

	data := struct {
		Name     string
		OldEmail string
		NewEmail string
		OTP      string
		LogoURL  string
	}{
		Name:     name,
		OldEmail: oldEmail,
		NewEmail: newEmail,
		OTP:      otp,
		LogoURL:  logoURL,
	}

	var content bytes.Buffer
	if err := tmpl.Execute(&content, data); err != nil {
		return buildSimpleEmailUpdateOTPEmail(name, oldEmail, newEmail, otp)
	}

	return content.String()
}

// buildSimpleEmailUpdateOTPEmail creates a simple HTML email for email update verification
func buildSimpleEmailUpdateOTPEmail(name, oldEmail, newEmail, otp string) string {
	return `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Verifikasi Perubahan Email</title>
		</head>
		<body>
			<h2>Halo, ` + name + `!</h2>
			<p>Kami menerima permintaan untuk mengubah alamat email Anda dari <strong>` + oldEmail + `</strong> menjadi <strong>` + newEmail + `</strong>.</p>
			<p>Berikut adalah kode verifikasi untuk konfirmasi perubahan email:</p>
			<h1 style="font-size: 32px; letter-spacing: 5px; background-color: #f5f5f5; padding: 10px; text-align: center;">` + otp + `</h1>
			<p>Kode ini akan kadaluarsa dalam 15 menit.</p>
			<p>Jika Anda tidak meminta perubahan email ini, silakan abaikan email ini dan hubungi administrator.</p>
			<p>Terima kasih,</p>
			<p>Tim LecSens</p>
		</body>
		</html>
	`
}

// BuildEmailChangeNotificationEmail creates the HTML content for email change notification to old email
func BuildEmailChangeNotificationEmail(name, oldEmail, newEmail string) string {
	templatePath := filepath.Join("utils", "email_template", "base-template-dark.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		// Fallback to simple template if file cannot be parsed
		return buildSimpleEmailChangeNotificationEmail(name, oldEmail, newEmail)
	}

	// Get logo URL from environment variable or use default
	logoURL := os.Getenv("EMAIL_LOGO_URL")
	if logoURL == "" {
		// Always use the app URL for logo reference
		logoURL = os.Getenv("APP_URL") + "/static/logo.png"
	}

	data := struct {
		Name     string
		OldEmail string
		NewEmail string
		LogoURL  string
	}{
		Name:     name,
		OldEmail: oldEmail,
		NewEmail: newEmail,
		LogoURL:  logoURL,
	}

	var content bytes.Buffer
	if err := tmpl.Execute(&content, data); err != nil {
		return buildSimpleEmailChangeNotificationEmail(name, oldEmail, newEmail)
	}

	return content.String()
}

// buildSimpleEmailChangeNotificationEmail creates a simple HTML email for notification to old email
func buildSimpleEmailChangeNotificationEmail(name, oldEmail, newEmail string) string {
	currentTime := time.Now().Format("02 Jan 2006, 15:04:05")
	return `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Pemberitahuan Perubahan Email</title>
		</head>
		<body>
			<h2>Halo, ` + name + `!</h2>
			<p>Kami ingin memberitahu Anda bahwa alamat email akun Anda telah berhasil diubah.</p>
			<p>Detail perubahan:</p>
			<ul>
				<li>Email lama: <strong>` + oldEmail + `</strong></li>
				<li>Email baru: <strong>` + newEmail + `</strong></li>
				<li>Waktu perubahan: <strong>` + currentTime + `</strong></li>
			</ul>
			<p>Jika Anda tidak melakukan perubahan ini, segera hubungi administrator kami.</p>
			<p>Terima kasih,</p>
			<p>Tim LecSens</p>
		</body>
		</html>
	`
}

// NewEmailSender creates a new EmailSender using environment variables
func NewEmailSender() EmailSender {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_AUTH_EMAIL")
	password := os.Getenv("SMTP_AUTH_PASSWORD")
	senderName := os.Getenv("SMTP_SENDER_NAME")

	// Convert port string to integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		// Default to port 587 if conversion fails
		port = 587
	}

	return &SMTPEmailSender{
		config: struct {
			Host     string
			Port     int
			Username string
			Password string
			Sender   string
		}{
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
			Sender:   fmt.Sprintf("%s <%s>", senderName, username),
		},
	}
}

// SaveEmailLogo saves the logo file to the assets directory
func SaveEmailLogo(logoData []byte) error {
	// Create assets directory if it doesn't exist
	assetsDir := filepath.Join("utils", "email_template", "assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		return fmt.Errorf("failed to create assets directory: %v", err)
	}

	// Save logo file
	logoPath := filepath.Join(assetsDir, "logo.png")
	if err := os.WriteFile(logoPath, logoData, 0644); err != nil {
		return fmt.Errorf("failed to save logo file: %v", err)
	}

	return nil
}

// GetEmailLogoPath returns the path to the logo file
func GetEmailLogoPath() string {
	return filepath.Join("utils", "email_template", "assets", "logo.png")
}
