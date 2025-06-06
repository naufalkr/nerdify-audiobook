package config

// EmailConfig holds the email configuration
type EmailConfig struct {
	Host         string
	Port         int
	SenderName   string
	AuthEmail    string
	AuthPassword string
}

// NewEmailConfig creates a new EmailConfig from the general Config
func NewEmailConfig() (*EmailConfig, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	return &EmailConfig{
		Host:         config.SMTPHost,
		Port:         config.SMTPPort,
		SenderName:   config.SMTPSenderName,
		AuthEmail:    config.SMTPAuthEmail,
		AuthPassword: config.SMTPAuthPassword,
	}, nil
}
