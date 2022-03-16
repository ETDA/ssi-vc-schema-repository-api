package home

import (
	core "ssi-gitlab.teda.th/ssi/core"
	"net/http"
)

type HomeController struct{}

func (hc *HomeController) Get(c core.IHTTPContext) error {
	return c.JSON(http.StatusOK, core.Map{
		"message": "foobar",
	})
}
