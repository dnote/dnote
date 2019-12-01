package mailer

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

// ErrSMTPNotConfigured is an error indicating that SMTP is not configured
var ErrSMTPNotConfigured = errors.New("SMTP is not configured")

// Backend is an interface for sending emails.
type Backend interface {
	Queue(subject, from string, to []string, contentType, body string) error
}

// SimpleBackendImplementation is an implementation of the Backend
// that sends an email without queueing.
type SimpleBackendImplementation struct {
}

type dialerParams struct {
	Host     string
	Port     int
	Username string
	Password string
}

func validateSMTPConfig() bool {
	port := os.Getenv("SmtpPort")
	host := os.Getenv("SmtpHost")
	username := os.Getenv("SmtpUsername")
	password := os.Getenv("SmtpPassword")

	return port != "" && host != "" && username != "" && password != ""
}

func getSMTPParams() (*dialerParams, error) {
	portEnv := os.Getenv("SmtpPort")
	hostEnv := os.Getenv("SmtpHost")
	usernameEnv := os.Getenv("SmtpUsername")
	passwordEnv := os.Getenv("SmtpPassword")

	if portEnv != "" && hostEnv != "" && usernameEnv != "" && passwordEnv != "" {
		return nil, ErrSMTPNotConfigured
	}

	port, err := strconv.Atoi(portEnv)
	if err != nil {
		return nil, errors.Wrap(err, "parsing SMTP port")
	}

	p := &dialerParams{
		Host:     hostEnv,
		Port:     port,
		Username: usernameEnv,
		Password: passwordEnv,
	}

	return p, nil
}

// Queue is an implementation of Backend.Queue.
func (b *SimpleBackendImplementation) Queue(subject, from string, to []string, contentType, body string) error {
	// If not production, never actually send an email
	if os.Getenv("GO_ENV") != "PRODUCTION" {
		log.Println("Not sending email because Dnote is not running in a production environment.")
		log.Printf("Subject: %s, to: %s, from: %s", subject, to, from)
		fmt.Println(body)
		return nil
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody(contentType, body)
	// m.SetBody("text/html", body)

	p, err := getSMTPParams()
	if err != nil {
		return errors.Wrap(err, "getting dialer params")
	}

	d := gomail.NewPlainDialer(p.Host, p.Port, p.Username, p.Password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
