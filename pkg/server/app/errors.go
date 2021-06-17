package app

type appError string

func (e appError) Error() string {
	return string(e)
}

func (e appError) Public() string {
	return string(e)
}

var (
	// ErrNotFound an error that indicates that the given resource is not found
	ErrNotFound appError = "not found"
	// ErrLoginInvalid is an error for invalid login
	ErrLoginInvalid appError = "Wrong email and password combination"

	// ErrDuplicateEmail is an error for duplicate email
	ErrDuplicateEmail appError = "duplicate email"
	// ErrEmailRequired is an error for missing email
	ErrEmailRequired appError = "Please enter an email"
	// ErrPasswordRequired is an error for missing email
	ErrPasswordRequired appError = "Please enter a password"
	// ErrPasswordTooShort is an error for short password
	ErrPasswordTooShort appError = "password should be longer than 8 characters"
	// ErrPasswordConfirmationMismatch is an error for password ans password confirmation not matching
	ErrPasswordConfirmationMismatch appError = "password confirmation does not match password"

	// ErrLoginRequired is an error for not authenticated
	ErrLoginRequired appError = "login required"

	// ErrBookUUIDRequired is an error for note missing book uuid
	ErrBookUUIDRequired appError = "book uuid required"
	// ErrBookNameRequired is an error for note missing book name
	ErrBookNameRequired appError = "book name required"
	// ErrDuplicateBook is an error for duplicate book
	ErrDuplicateBook appError = "duplicate book exists"

	// ErrEmptyUpdate is an error for empty update params
	ErrEmptyUpdate appError = "update is empty"

	// ErrInvalidUUID is an error for invalid uuid
	ErrInvalidUUID appError = "invalid uuid"

	// ErrInvalidSMTPConfig is an error for invalid SMTP configuration
	ErrInvalidSMTPConfig appError = "SMTP is not configured"

	// ErrInvalidToken is an error for invalid token
	ErrInvalidToken appError = "invalid token"

	// ErrPasswordResetTokenExpired is an error for expired password reset token
	ErrPasswordResetTokenExpired appError = "this link has been expired. Please request a new password reset link."
)
