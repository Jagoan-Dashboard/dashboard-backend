package validation

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	"building-report-backend/internal/domain/constants"
)

var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidUsername = errors.New("invalid username format")
	ErrPasswordTooWeak = errors.New("password too weak")
	ErrTextTooLong     = errors.New("text exceeds maximum length")
	ErrInvalidURL      = errors.New("invalid URL format")
	ErrCoordinateRange = errors.New("coordinates out of valid range")
)

// Email validation using regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Username validation (alphanumeric and underscore, 3-50 chars)
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,50}$`)

// URL validation
var urlRegex = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)

// ValidateEmail checks if email format is valid
func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

// ValidateUsername checks if username format is valid
func ValidateUsername(username string) error {
	if !usernameRegex.MatchString(username) {
		return ErrInvalidUsername
	}
	return nil
}

// ValidatePassword checks password strength
// Requirements: min 8 chars, at least 1 upper, 1 lower, 1 number
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooWeak
	}

	var hasUpper, hasLower, hasNumber bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return ErrPasswordTooWeak
	}

	return nil
}

// ValidateTextLength checks if text exceeds maximum length
func ValidateTextLength(text string, maxLength int) error {
	if len(strings.TrimSpace(text)) > maxLength {
		return ErrTextTooLong
	}
	return nil
}

// ValidateURL checks if URL format is valid
func ValidateURL(url string) error {
	if url == "" {
		return nil // Empty URL is allowed
	}

	if !urlRegex.MatchString(url) {
		return ErrInvalidURL
	}
	return nil
}

// ValidateRequired checks if required field is not empty
func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(fieldName + " is required")
	}
	return nil
}

// ValidateCoordinates checks if latitude and longitude are in valid ranges
func ValidateCoordinates(lat, lng float64) error {
	if lat < -90 || lat > 90 {
		return ErrCoordinateRange
	}
	if lng < -180 || lng > 180 {
		return ErrCoordinateRange
	}
	return nil
}

// ValidatePositiveNumber checks if a number is positive
func ValidatePositiveNumber(value float64, fieldName string) error {
	if value < 0 {
		return errors.New(fieldName + " must be positive")
	}
	return nil
}

// ValidateInRange checks if a number is within specified range
func ValidateInRange(value, min, max float64, fieldName string) error {
	if value < min || value > max {
		return errors.New(fieldName + " must be between " +
			string(rune(min)) + " and " + string(rune(max)))
	}
	return nil
}

// ValidatePageSize validates pagination page size
func ValidatePageSize(size int) int {
	if size <= 0 {
		return constants.DefaultPageSize
	}
	if size > constants.MaxPageSize {
		return constants.MaxPageSize
	}
	return size
}

// ValidatePage validates pagination page number
func ValidatePage(page int) int {
	if page <= 0 {
		return 1
	}
	return page
}