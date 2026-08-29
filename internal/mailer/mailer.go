package mailer

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

type Mailer struct {
	client *resend.Client
	from   string
}

func New(apiKey, from string) *Mailer {
	return &Mailer{
		client: resend.NewClient(apiKey),
		from:   from,
	}
}

func (m *Mailer) SendOTP(toEmail, name, otp string) error {
	html := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2>Verify your email</h2>
			<p>Hi %s,</p>
			<p>Your verification code is:</p>
			<div style="
				background: #f4f4f4;
				padding: 20px;
				text-align: center;
				font-size: 36px;
				font-weight: bold;
				letter-spacing: 8px;
				border-radius: 8px;
				margin: 20px 0;
			">%s</div>
			<p>This code expires in <strong>10 minutes</strong>.</p>
			<p>If you didn't request this, ignore this email.</p>
		</div>
	`, name, otp)

	params := &resend.SendEmailRequest{
		From:    m.from,
		To:      []string{toEmail},
		Subject: "Your verification code",
		Html:    html,
	}

	_, err := m.client.Emails.Send(params)
	return err
}

func (m *Mailer) SendResetOTP(toEmail, name, otp string) error {
	html := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2>Reset your password</h2>
			<p>Hi %s,</p>
			<p>Your password reset code is:</p>
			<div style="
				background: #f4f4f4;
				padding: 20px;
				text-align: center;
				font-size: 36px;
				font-weight: bold;
				letter-spacing: 8px;
				border-radius: 8px;
				margin: 20px 0;
			">%s</div>
			<p>This code expires in <strong>10 minutes</strong>.</p>
			<p>If you didn't request a password reset, ignore this email.</p>
		</div>
	`, name, otp)

	params := &resend.SendEmailRequest{
		From:    m.from,
		To:      []string{toEmail},
		Subject: "Reset your password",
		Html:    html,
	}

	_, err := m.client.Emails.Send(params)
	return err
}

func (m *Mailer) SendResendOTP(toEmail, name, otp string) error {
	html := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
			<h2>Your new verification code</h2>
			<p>Hi %s,</p>
			<p>Here is your new verification code:</p>
			<div style="
				background: #f4f4f4;
				padding: 20px;
				text-align: center;
				font-size: 36px;
				font-weight: bold;
				letter-spacing: 8px;
				border-radius: 8px;
				margin: 20px 0;
			">%s</div>
			<p>This code expires in <strong>10 minutes</strong>.</p>
			<p>If you didn't request this, ignore this email.</p>
		</div>
	`, name, otp)

	params := &resend.SendEmailRequest{
		From:    m.from,
		To:      []string{toEmail},
		Subject: "Your new verification code",
		Html:    html,
	}

	_, err := m.client.Emails.Send(params)
	return err
}
