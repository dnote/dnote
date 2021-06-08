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
	ErrLoginInvalid appError = "invalid login"

	// ErrDuplicateEmail is an error for duplicate email
	ErrDuplicateEmail appError = "duplicate email"
	// ErrEmailRequired is an error for missing email
	ErrEmailRequired appError = "missing email"
	// ErrPasswordTooShort is an error for short password
	ErrPasswordTooShort appError = "password should be longer than 8 characters"

	// ErrLoginRequired is an error for not authenticated
	ErrLoginRequired appError = "login required"

	// ErrBookUUIDRequired is an error for note missing book uuid
	ErrBookUUIDRequired appError = "book uuid required"

	// ErrEmptyUpdate is an error for empty update params
	ErrEmptyUpdate appError = "update is empty"

	// ErrInvalidUUID is an error for invalid uuid
	ErrInvalidUUID appError = "invalid uuid"
)
