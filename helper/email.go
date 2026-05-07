package helper

import (
	"Zhooze/config"
	"fmt"
	"net/smtp"
)

func SendEmail(to string, subject string, body string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Only send if configured
	if cfg.SMTPHost == "" || cfg.SMTPUser == "" {
		fmt.Println("SMTP not configured, skipping email.")
		return nil
	}

	from := cfg.SMTPUser
	password := cfg.SMTPPass
	smtpHost := cfg.SMTPHost
	smtpPort := cfg.SMTPPort

	auth := smtp.PlainAuth("", from, password, smtpHost)

	// RFC 822 format
	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body))

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	err = smtp.SendMail(addr, auth, from, []string{to}, message)
	if err != nil {
		return err
	}
	return nil
}

func SendWelcomeEmail(userEmail string, userName string) error {
	fmt.Printf("Entered SendWelcomeEmail for: %s\n", userEmail)
	subject := "Welcome to Soxenly!"
	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
			<h2 style="color: #1B4332;">Welcome to Soxenly!</h2>
			<p>Hello %s,</p>
			<p>Thank you for joining Soxenly! We are excited to have you with us.</p>
			<p>Explore our sustainable collection and enjoy premium comfort.</p>
			<div style="margin-top: 30px;">
				<a href="https://soxenly.vercel.app/" style="background-color: #1B4332; color: white; padding: 10px 20px; text-decoration: none;">Visit Store</a>
			</div>
			<p style="margin-top: 40px; font-size: 10px; color: #888;">© 2026 Soxenly</p>
		</div>
	`, userName)

	return SendEmail(userEmail, subject, body)
}

func SendOrderConfirmationEmail(userEmail string, userName string, orderID int, totalAmount float64) error {
	subject := fmt.Sprintf("Order Confirmation #%d - Soxenly", orderID)
	body := fmt.Sprintf(`
		<div style="font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; color: #1B4332;">
			<div style="text-align: center; margin-bottom: 30px;">
				<h1 style="color: #1B4332; margin: 0;">SOXENLY</h1>
				<p style="text-transform: uppercase; letter-spacing: 2px; font-size: 10px; color: #40916C;">Order Confirmed</p>
			</div>
			<p>Hello %s,</p>
			<p>Thank you for your purchase! Your order <strong>#%d</strong> has been successfully placed and is being prepared for shipment.</p>
			<div style="background-color: #D8F3DC; padding: 20px; margin: 20px 0;">
				<p style="margin: 0; font-size: 14px;"><strong>Order Total:</strong> ₹%.2f</p>
				<p style="margin: 5px 0 0 0; font-size: 12px; opacity: 0.8;">You will receive another email once your items have shipped.</p>
			</div>
			<p>Thank you for choosing conscious comfort.</p>
			<div style="text-align: center; margin: 40px 0;">
				<a href="https://soxenly.vercel.app/orders" style="background-color: #1B4332; color: #D8F3DC; padding: 15px 30px; text-decoration: none; font-weight: bold; text-transform: uppercase; font-size: 12px; letter-spacing: 1px;">View Order Status</a>
			</div>
			<p style="font-size: 12px; color: #1B4332; opacity: 0.6; margin-top: 50px; border-top: 1px solid #D8F3DC; pt-20">
				© 2026 SOXENLY. All rights reserved.
			</p>
		</div>
	`, userName, orderID, totalAmount)

	return SendEmail(userEmail, subject, body)
}
