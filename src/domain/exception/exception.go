package exception

import (
	"errors"
	"fmt"
)

var DomainError = errors.New("")
var NotFoundError = fmt.Errorf("not found%w", DomainError)

func New(text string) error {
	return fmt.Errorf("%s%w", text, DomainError)
}

func Errorf(format string, a ...any) error {
	return New(fmt.Sprintf(format, a...))
}

func IsDomainException(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, DomainError)
}

func IsNotFoundException(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, NotFoundError)
}
