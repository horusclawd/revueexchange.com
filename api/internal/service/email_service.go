package service

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// EmailService handles sending emails
type EmailService struct {
	fromEmail string
	fromName  string
	apiKey    string
}

// NewEmailService creates a new email service
func NewEmailService(apiKey, fromEmail, fromName string) *EmailService {
	return &EmailService{
		fromEmail: fromEmail,
		fromName:  fromName,
		apiKey:    apiKey,
	}
}

// SendEmail sends an email
func (s *EmailService) SendEmail(to, subject, htmlContent string) error {
	if s.apiKey == "" {
		// No API key configured - skip sending
		return nil
	}

	from := mail.NewEmail(s.fromName, s.fromEmail)
	toEmail := mail.NewEmail("", to)

	message := mail.NewSingleEmail(from, subject, toEmail, "", htmlContent)
	client := sendgrid.NewSendClient(s.apiKey)

	_, err := client.Send(message)
	return err
}

// SendWelcomeEmail sends a welcome email
func (s *EmailService) SendWelcomeEmail(to, username string) error {
	subject := "Welcome to RevUExchange!"
	htmlContent := fmt.Sprintf(`
		<h1>Welcome to RevUExchange, %s!</h1>
		<p>Thanks for joining our community of authors and reviewers.</p>
		<p>Start by:</p>
		<ul>
			<li>Adding your products (books, courses, newsletters)</li>
			<li>Browsing bounties to review</			<li>Earning points by writing quality reviews</li>
		</ul>
		<p>We look forward to seeing your reviews!</p>
	`, username)
	return s.SendEmail(to, subject, htmlContent)
}

// SendReviewReceivedEmail notifies author of new review
func (s *EmailService) SendReviewReceivedEmail(to, authorName, productTitle string) error {
	subject := "New Review on Your Product!"
	htmlContent := fmt.Sprintf(`
		<h1>Great news, %s!</h1>
		<p>You've received a new review on "%s".</p>
		<p>Log in to check it out and respond to the reviewer.</p>
	`, authorName, productTitle)
	return s.SendEmail(to, subject, htmlContent)
}

// SendBountyCompletedEmail notifies author bounty is complete
func (s *EmailService) SendBountyCompletedEmail(to, authorName, productTitle string) error {
	subject := "Your Bounty is Complete!"
	htmlContent := fmt.Sprintf(`
		<h1>Your bounty is complete, %s!</h1>
		<p>The review for "%s" has been submitted.</p>
		<p>Log in to see the review and award points to the reviewer.</p>
	`, authorName, productTitle)
	return s.SendEmail(to, subject, htmlContent)
}

// SendPointsAwardedEmail notifies user of points earned
func (s *EmailService) SendPointsAwardedEmail(to, username string, points int) error {
	subject := "You Earned Points!"
	htmlContent := fmt.Sprintf(`
		<h1>Congratulations, %s!</h1>
		<p>You've earned %d points for your review!</p>
		<p>Use your points to create bounties for your own products.</p>
	`, username, points)
	return s.SendEmail(to, subject, htmlContent)
}

// InitEmailService initializes email service from config
func InitEmailService(apiKey, fromEmail, fromName string) *EmailService {
	if apiKey == "" {
		return nil
	}
	return NewEmailService(apiKey, fromEmail, fromName)
}
