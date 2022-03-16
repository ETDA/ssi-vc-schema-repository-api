package home

import (
	"github.com/labstack/echo/v4"
	core "ssi-gitlab.teda.th/ssi/core"
)

func NewHomeHandler(r *echo.Echo) {

	home := &HomeController{}
	r.GET("/", core.WithHTTPContext(home.Get))
}
