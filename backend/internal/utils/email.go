package utils

import (
	"backend/internal/config"
	"fmt"
	"log"
	"net/smtp"
)

func SendEmail(to string, subject string, body string) error {
	cfg := config.Load()

	log.Printf("Tentative d'envoi d'email à %s (SMTP: %s)", to, cfg.SMTPHost)

	// If no SMTP config, just log it (useful for local dev)
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

	// Mailpit doesn't need auth, but SendMail requires auth object or nil
	var auth smtp.Auth
	if cfg.SMTPUser != "" {
		auth = smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)
	}

	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)

	err := smtp.SendMail(addr, auth, cfg.FromEmail, []string{to}, msg)
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
