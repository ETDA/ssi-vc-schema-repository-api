package token

import (
	"github.com/labstack/echo/v4"
	"gitlab.finema.co/finema/etda/vc-schema-api/middlewares"
	core "ssi-gitlab.teda.th/ssi/core"
)

func NewSchemaHandler(r *echo.Echo) {

	token := &TokenController{}

	r.GET("/schemas/tokens", core.WithHTTPContext(token.Pagination), middlewares.IsAuth, middlewares.IsAdmin)
	r.GET("/schemas/tokens/:id", core.WithHTTPContext(token.Find), middlewares.IsAuth, middlewares.IsAdmin)
	r.POST("/schemas/tokens", core.WithHTTPContext(token.Create), middlewares.IsAuth, middlewares.IsAdmin)
	r.PUT("/schemas/tokens/:id", core.WithHTTPContext(token.Update), middlewares.IsAuth, middlewares.IsAdmin)
	r.DELETE("/schemas/tokens/:id", core.WithHTTPContext(token.Delete), middlewares.IsAuth, middlewares.IsAdmin)
	r.GET("/schemas/tokens/me", core.WithHTTPContext(token.Me), middlewares.IsAuth)

}
