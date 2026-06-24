package utils

import (
	"backend/internal/config"
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

func SendEmail(to string, subject string, body string) error {
	cfg := config.Load()

	log.Printf("Tentative d'envoi d'email à %s (SMTP: %s, Port: %s)", to, cfg.SMTPHost, cfg.SMTPPort)

	if cfg.SMTPHost == "" {
		log.Printf("SIMULATION EMAIL to %s\nSubject: %s\nBody: %s\n", to, subject, body)
		return nil
	}

	msg := []byte("From: " + cfg.FromEmail + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)

	var err error
	if cfg.SMTPPort == "465" {
		// SSL/TLS implicit connection for port 465
		tlsconfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         cfg.SMTPHost,
		}
		conn, dialErr := tls.Dial("tcp", addr, tlsconfig)
		if dialErr != nil {
			log.Printf("ERREUR SMTP TLS Dial: %v", dialErr)
			return fmt.Errorf("erreur lors de l'envoi de l'email (TLS dial): %w", dialErr)
		}
		defer conn.Close()

		c, clientErr := smtp.NewClient(conn, cfg.SMTPHost)
		if clientErr != nil {
			log.Printf("ERREUR SMTP Client creation: %v", clientErr)
			return fmt.Errorf("erreur lors de l'envoi de l'email (SMTP client): %w", clientErr)
		}
		defer c.Quit()

		if cfg.SMTPUser != "" {
			auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)
			if authErr := c.Auth(auth); authErr != nil {
				log.Printf("ERREUR SMTP Auth: %v", authErr)
				return fmt.Errorf("erreur lors de l'envoi de l'email (Auth): %w", authErr)
			}
		}

		if err = c.Mail(cfg.FromEmail); err != nil {
			return err
		}
		if err = c.Rcpt(to); err != nil {
			return err
		}
		w, dataErr := c.Data()
		if dataErr != nil {
			return dataErr
		}
		_, err = w.Write(msg)
		if err != nil {
			return err
		}
		err = w.Close()
		if err != nil {
			return err
		}
	} else {
		// Standard smtp.SendMail for STARTTLS (usually port 587) or unencrypted
		var auth smtp.Auth
		if cfg.SMTPUser != "" {
			auth = smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)
		}
		err = smtp.SendMail(addr, auth, cfg.FromEmail, []string{to}, msg)
	}
	if err != nil {
		log.Printf("ERREUR SMTP: %v", err)
		return fmt.Errorf("erreur lors de l'envoi de l'email: %w", err)
	}

	log.Printf("Email envoyé avec succès à %s", to)
	return nil
}

func SendConfirmationEmail(to string, username string) error {
	subject := "Bienvenue chez UpCycleConnect !"
	body := fmt.Sprintf(`
		<h1>Bonjour %s !</h1>
		<p>Merci de vous être inscrit sur <strong>UpCycleConnect</strong>.</p>
		<p>Votre compte a bien été créé. Vous pouvez dès à présent proposer des objets, participer à des ateliers ou rejoindre des collectes.</p>
		<p>À très bientôt sur la plateforme !</p>
		<br/>
		<p>L'équipe UpCycleConnect</p>
	`, username)

	return SendEmail(to, subject, body)
}
