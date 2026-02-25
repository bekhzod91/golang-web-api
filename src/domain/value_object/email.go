package value_object

import (
	"github.com/myproject/api/domain/exception"
	"regexp"
	"strings"
)

type Email string

// ParseEmail validates and returns a new Email object
func ParseEmail(value string) (Email, error) {
	if !isValidEmail(value) {
		return "", exception.New("invalid email address format")
	}
	return Email(strings.ToLower(strings.TrimSpace(value))), nil
}

// isValidEmail performs a basic email format validation using a regex
func isValidEmail(email string) bool {
	// Basic email regex for demonstration; can be made more comprehensive
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// GetValue returns the underlying string representation of the Email
func (e Email) String() string {
	return string(e)
}
