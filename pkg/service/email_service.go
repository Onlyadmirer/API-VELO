package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/resend/resend-go/v3"
)

type EmailService interface {
	SendVerificationEmail(toEmail string, name string, verifyToken string) error
}

type resendEmailService struct {
	client *resend.Client
}

func NewEmailService() EmailService {
	apiKey := os.Getenv("RESEND_KEY")
	client := resend.NewClient(apiKey)
	return &resendEmailService{
		client: client,
	}
}

func (s *resendEmailService) SendVerificationEmail(toEmail string, name string, verifyToken string) error {

	verifyURL := fmt.Sprintf("http://localhost:3000/verify-email?token=%s", verifyToken)

	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; line-height: 1.5; color: #333;">
			<h2>Halo %s! Selamat datang di VELO.</h2>
			<p>Terima kasih telah mendaftar. Tinggal satu langkah lagi untuk mengaktifkan akunmu.</p>
			<p>Silakan klik tombol di bawah ini untuk memverifikasi emailmu:</p>
			<a href="%s" style="display: inline-block; padding: 10px 20px; margin-top: 10px; background-color: #000000; color: #ffffff; text-decoration: none; border-radius: 5px; font-weight: bold;">Verifikasi Akun</a>
			<p style="margin-top: 20px; font-size: 12px; color: #777;">Jika kamu tidak mendaftar di VELO, abaikan email ini.</p>
		</div>
	`, name, verifyURL)

	params := &resend.SendEmailRequest{
		From:    "Akmaldev@akmaldev.my.id",
		To:      []string{toEmail},
		Subject: "Verifikasi Akun VELO",
		Html:    htmlBody,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sent, err := s.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("gagal kirim email ke: %s: %v", toEmail, err)
	}

	log.Printf("Email verifikasi berhasil dikirim ke: %s (Resend ID: %s)", toEmail, sent.Id)

	return nil

}
