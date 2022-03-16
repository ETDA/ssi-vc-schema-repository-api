package token

import (
	"gitlab.finema.co/finema/etda/vc-schema-api/consts"
	"net/http"

	"gitlab.finema.co/finema/etda/vc-schema-api/requests"
	"gitlab.finema.co/finema/etda/vc-schema-api/services"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type TokenController struct {
}

func (tc *TokenController) Find(c core.IHTTPContext) error {
	s := services.NewTokenService(c)
	token, ierr := s.FindByID(c.Param("id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, token)
}

func (tc TokenController) Create(c core.IHTTPContext) error {
	input := &requests.TokenCreate{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	s := services.NewTokenService(c)

	token, ierr := s.Create(&services.CreateTokenPayload{
		Name: utils.GetString(input.Name),
		Role: utils.GetString(input.Role),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusCreated, token)
}

func (tc TokenController) Update(c core.IHTTPContext) error {
	input := &requests.TokenUpdate{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	s := services.NewTokenService(c)

	token, ierr := s.Update(&services.UpdateTokenPayload{
		ID:   c.Param("id"),
		Name: utils.GetString(input.Name),
		Role: utils.GetString(input.Role),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, token)
}

func (tc TokenController) Pagination(c core.IHTTPContext) error {
	s := services.NewTokenService(c)
	tokens, pageResponse, ierr := s.Pagination(c.GetPageOptionsWithOptions(&core.PageOptionsOptions{}))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, core.NewPagination(tokens, pageResponse))
}

func (tc *TokenController) Delete(c core.IHTTPContext) error {
	s := services.NewTokenService(c)
	ierr := s.Delete(c.Param("id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.NoContent(http.StatusNoContent)
}

func (tc *TokenController) Me(c core.IHTTPContext) error {
	tokenID := c.Get(consts.ContextKeyTokenID).(string)
	s := services.NewTokenService(c)
	token, ierr := s.FindByID(tokenID)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, token)
}
