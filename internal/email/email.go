package email

import (
	gomail "gopkg.in/mail.v2"
)

type EmailParams struct {
	To      string
	Subject string
	Body    string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func SendEmail(cfg SMTPConfig, params EmailParams) error {
	message := gomail.NewMessage()
	message.SetHeader("From", cfg.Username)
	message.SetHeader("To", params.To)
	message.SetHeader("Subject", params.Subject)
	message.SetBody("text/html", params.Body)

	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	if err := dialer.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
