package mailer

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendPasswordResetEmail(to, resetLink string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")

	auth := smtp.PlainAuth("", user, pass, host)

	subject := "Recuperacao de senha - Gestio"
	body := fmt.Sprintf(
		"Para redefinir sua senha, clique no link abaixo:\n\n%s\n\nSe voce nao solicitou isso, ignore este email.",
		resetLink,
	)

	message := []byte("Subject: " + subject + "\r\n" +
		"From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body)

	return smtp.SendMail(host+":"+port, auth, user, []string{to}, message)
}
