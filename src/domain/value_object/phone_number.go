package value_object

import (
	"regexp"
	"strings"

	"github.com/myproject/api/domain/exception"
)

type PhoneNumber string

// ParsePhoneNumber validates and returns a new PhoneNumber object
func ParsePhoneNumber(value string) (PhoneNumber, error) {
	if !isValidPhoneNumber(value) {
		return "", exception.New("invalid phone number format")
	}
	return PhoneNumber(strings.TrimSpace(value)), nil
}

// isValidPhoneNumber performs a basic phone number format validation using a regex
func isValidPhoneNumber(phoneNumber string) bool {
	// Basic phone number regex for demonstration; can be made more comprehensive
	const phoneNumberRegex = `^[1-9][0-9]{7,14}$`
	re := regexp.MustCompile(phoneNumberRegex)
	return re.MatchString(phoneNumber)
}

// GetValue returns the underlying string representation of the PhoneNumber
func (e PhoneNumber) String() string {
	return string(e)
}
