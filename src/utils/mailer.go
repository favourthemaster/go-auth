package utils

import "fmt"

type Mailer interface {
	SendMail(to, content string) error
	SendVerificationMail(to, verificationLink string) error
	SendPasswordResetMail(to, passwordResetLink string) error
	SendPasswordChangeMail(to, passwordChangeLink string) error
}

type mailer struct {
}

func NewMailer() Mailer {
	return &mailer{}
}

func (m mailer) SendMail(to, content string) error {
	fmt.Printf("TO: %s, Content: %s\n", to, content)
	return nil
}

func (m mailer) SendVerificationMail(to, verificationLink string) error {
	fmt.Printf("TO: %s, Content: %s\n", to, verificationLink)
	return nil
}

func (m mailer) SendPasswordResetMail(to, passwordResetLink string) error {
	fmt.Printf("TO: %s, Content: %s\n", to, passwordResetLink)
	return nil
}

func (m mailer) SendPasswordChangeMail(to, passwordChangeLink string) error {
	fmt.Printf("TO: %s, Content: %s\n", to, passwordChangeLink)
	return nil
}
