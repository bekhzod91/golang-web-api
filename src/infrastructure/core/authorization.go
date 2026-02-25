package core

import (
	gocontenxt "context"
	"errors"
	"net/http"
	"strings"

	"github.com/myproject/api/domain/entity"
	"github.com/myproject/api/domain/value_object"
)

var ContextKeyUser = "_core/user"

func Authorization(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c := NewContext(w, r)

		header := r.Header.Get("Authorization")
		if header == "" {
			err := errors.New("token not provided")
			c.JSON(http.StatusUnauthorized, M{"message": err.Error()})
			return
		}

		tokenString, err := parseAuthorizationHeader(header)
		if err != nil {
			c.JSON(http.StatusUnauthorized, M{"message": err.Error()})
			return
		}

		token, err := value_object.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, M{"message": err.Error()})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, M{"message": err.Error()})
			return
		}

		var user *entity.User
		user, err = c.Storage().Token().GetUserByToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, M{"message": err.Error()})
			return
		}

		ctx := r.Context()
		ctx = gocontenxt.WithValue(ctx, ContextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func parseAuthorizationHeader(header string) (string, error) {
	value := strings.Split(strings.TrimSpace(header), " ")

	if len(value) != 2 {
		return "", errors.New("invalid token")
	}

	if value[0] != "Bearer" {
		return "", errors.New("invalid token")
	}

	return value[1], nil
}

func UserFromContext(ctx IContext) *entity.User {
	user, ok := ctx.Request().Context().Value(ContextKeyUser).(*entity.User)
	if ok {
		return user
	}

	panic(errors.New("couldn't get web from context"))
}
