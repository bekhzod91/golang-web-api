package value_object

import (
	"github.com/hzmat24/api/domain/exception"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

const minLength = 8

type Password string

func NewPassword(password string) (Password, error) {
	if len(password) < minLength {
		return "", exception.New("password length should be at least 8 characters")
	}

	if strings.ToLower(password) == password {
		return "", exception.New("password should contain at least one uppercase character")
	}

	if strings.ToUpper(password) == password {
		return "", exception.New("password should contain at least one lowercase character")
	}

	hash, err := hashPassword(password)
	if err != nil {
		return "", err
	}

	return Password(hash), nil
}

func (p Password) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p), []byte(password))
	return err == nil
}

func (p Password) String() string {
	return string(p)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
