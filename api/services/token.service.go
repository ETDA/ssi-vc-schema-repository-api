package services

import (
	"errors"

	"gitlab.finema.co/finema/etda/vc-schema-api/consts"
	"gitlab.finema.co/finema/etda/vc-schema-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"gorm.io/gorm"
)

type CreateTokenPayload struct {
	Name string
	Role string
}

type UpdateTokenPayload struct {
	ID   string
	Name string
	Role string
}

type ITokenService interface {
	Create(payload *CreateTokenPayload) (*models.Token, core.IError)
	Update(payload *UpdateTokenPayload) (*models.Token, core.IError)
	FindByID(id string) (*models.Token, core.IError)
	FindAdminToken() (*models.Token, core.IError)
	FindByToken(token string) (*models.Token, core.IError)
	Pagination(pageOptions *core.PageOptions) ([]models.Token, *core.PageResponse, core.IError)
	Delete(id string) core.IError
}

type tokenService struct {
	ctx core.IContext
}

func NewTokenService(ctx core.IContext) ITokenService {
	return &tokenService{ctx: ctx}
}

func (s tokenService) Create(payload *CreateTokenPayload) (*models.Token, core.IError) {
	id := utils.GetUUID()
	createdAt := utils.GetCurrentDateTime()
	token := &models.Token{
		ID:        utils.GetUUID(),
		Name:      payload.Name,
		Token:     utils.NewSha256(id + createdAt.String()),
		Role:      payload.Role,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
	err := s.ctx.DB().Create(token).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return s.FindByID(token.ID)
}

func (s tokenService) Update(payload *UpdateTokenPayload) (*models.Token, core.IError) {
	token, ierr := s.FindByID(payload.ID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	token.Name = payload.Name
	token.Role = payload.Role
	token.UpdatedAt = utils.GetCurrentDateTime()

	err := s.ctx.DB().Updates(token).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return s.FindByID(token.ID)
}

func (s tokenService) FindByID(id string) (*models.Token, core.IError) {
	item := &models.Token{}

	err := s.ctx.DB().Where("id = ?", id).First(item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFound)
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return item, nil
}

func (s tokenService) FindByToken(token string) (*models.Token, core.IError) {
	item := &models.Token{}

	err := s.ctx.DB().Where("token = ?", token).First(item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFound)
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return item, nil
}

func (s tokenService) FindAdminToken() (*models.Token, core.IError) {
	item := &models.Token{}

	err := s.ctx.DB().Where("role = ?", consts.TokenAdminRole).First(item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFound)
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return item, nil
}

func (s tokenService) Pagination(pageOptions *core.PageOptions) ([]models.Token, *core.PageResponse, core.IError) {
	items := make([]models.Token, 0)

	if len(pageOptions.OrderBy) == 0 {
		pageOptions.OrderBy = []string{"created_at desc"}

	}

	pageRes, err := core.Paginate(s.ctx.DB(), &items, pageOptions)
	if err != nil {
		return nil, nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return items, pageRes, nil
}

func (s tokenService) Delete(id string) core.IError {
	token, ierr := s.FindByID(id)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	err := s.ctx.DB().Delete(token).Error
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}

	return nil
}
