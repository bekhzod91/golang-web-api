package value_object

import (
	"fmt"
	"github.com/myproject/api/domain/exception"
	"github.com/myproject/api/pkg/randutil"
	"strings"
)

var tokenPrefix = "auth"

type Token string

func NewToken(userId int64) (Token, error) {
	tokenSuffix, err := randutil.GenerateRandomStrings(30)
	if err != nil {
		return "", fmt.Errorf("failed creating token: %w", err)
	}
	return Token(fmt.Sprintf("%s:%d:%s", tokenPrefix, userId, tokenSuffix)), nil
}

func ParseToken(tokenString string) (Token, error) {
	if strings.HasPrefix(tokenString, tokenPrefix) {
		return Token(tokenString), nil
	}

	return "", exception.New(fmt.Sprintf("invalid token, token prefix should be %s", tokenPrefix))
}

func (t Token) String() string {
	return string(t)
}
