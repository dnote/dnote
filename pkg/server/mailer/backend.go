package mailer

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

// Backend is an interface for sending emails.
type Backend interface {
	Queue(subject, from string, to []string, body string) error
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

func getSMTPParams() (*dialerParams, error) {
	portStr := os.Getenv("SmtpPort")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.Wrap(err, "parsing SMTP port")
	}

	p := &dialerParams{
		Host:     os.Getenv("SmtpHost"),
		Port:     port,
		Username: os.Getenv("SmtpUsername"),
		Password: os.Getenv("SmtpPassword"),
	}

	return p, nil
}

// Queue is an implementation of Backend.Queue.
func (b *SimpleBackendImplementation) Queue(subject, from string, to []string, body string) error {
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
	m.SetBody("text/html", body)

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
