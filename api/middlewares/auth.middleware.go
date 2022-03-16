package middlewares

import (
	"github.com/labstack/echo/v4"
	"gitlab.finema.co/finema/etda/vc-schema-api/consts"
	"gitlab.finema.co/finema/etda/vc-schema-api/emsgs"
	"gitlab.finema.co/finema/etda/vc-schema-api/services"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"strings"
)

func IsAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(core.IHTTPContext)
		authentication := strings.TrimSpace(cc.Request().Header.Get("Authorization"))
		if authentication == "" {
			return c.JSON(emsgs.AuthTokenRequired.GetStatus(), emsgs.AuthTokenRequired.JSON())
		}
		splittedAuthentication := strings.Split(authentication, " ")

		if len(splittedAuthentication) != 2 {
			return c.JSON(emsgs.AuthTokenInvalid.GetStatus(), emsgs.AuthTokenInvalid.JSON())
		}
		prefix := splittedAuthentication[0]
		token := splittedAuthentication[1]
		if prefix != cc.ENV().String(consts.ENVTokenPrefix) {
			return c.JSON(emsgs.AuthTokenInvalid.GetStatus(), emsgs.AuthTokenInvalid.JSON())
		}

		service := services.NewTokenService(cc)
		t, ierr := service.FindByToken(token)
		if errmsgs.IsNotFoundError(ierr) {
			return c.JSON(emsgs.AuthTokenInvalid.GetStatus(), emsgs.AuthTokenInvalid.JSON())
		}
		if ierr != nil {
			return c.JSON(ierr.GetStatus(), ierr.JSON())
		}

		cc.Set(consts.ContextKeyToken, token)
		cc.Set(consts.ContextKeyTokenID, t.ID)

		return next(c)
	}
}
