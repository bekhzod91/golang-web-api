package command

import (
	"strings"

	"github.com/hzmat24/api/domain/value_object"
	"github.com/hzmat24/api/infrastructure/core"
)

func SignOut(ctx core.IContext) error {
	header := ctx.Request().Header.Get("Authorization")
	headerTokenParts := strings.Split(header, " ")

	token, err := value_object.ParseToken(headerTokenParts[1])
	if err != nil {
		return err
	}

	err = ctx.Storage().Token().RevokeUserToken(token)
	if err != nil {
		return err
	}

	return nil
}
