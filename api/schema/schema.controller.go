package schema

import (
	"encoding/json"
	"net/http"
	"strings"

	"gitlab.finema.co/finema/etda/vc-schema-api/consts"
	"gitlab.finema.co/finema/etda/vc-schema-api/emsgs"
	"gitlab.finema.co/finema/etda/vc-schema-api/helpers"
	"gitlab.finema.co/finema/etda/vc-schema-api/requests"
	"gitlab.finema.co/finema/etda/vc-schema-api/services"
	"gitlab.finema.co/finema/etda/vc-schema-api/views"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type SchemaController struct{}

func (sm SchemaController) Create(c core.IHTTPContext) error {
	input := &requests.VCSchemaCreate{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	s := services.NewSchemaService(c)

	schemaBody := json.RawMessage(utils.JSONToString(input.SchemaBody))
	createdBy := c.Get(consts.ContextKeyTokenID).(string)
	vcSchema, ierr := s.Create(&services.SchemaCreatePayload{
		SchemaName: utils.GetString(input.SchemaName),
		SchemaType: utils.GetString(input.SchemaType),
		SchemaBody: &schemaBody,
		CreatedBy:  createdBy,
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusCreated, vcSchema)
}

func (sm SchemaController) Upload(c core.IHTTPContext) error {
	formFile, err := c.FormFile("file")
	if err != nil {
		return c.JSON(emsgs.RequestFormError(err).GetStatus(), emsgs.RequestFormError(err).JSON())
	}

	if helpers.FileExtensionMatch(formFile.Filename, "zip") != nil {
		return c.JSON(emsgs.IsNotZipFile().GetStatus(), emsgs.IsNotZipFile().JSON())
	}

	zipJSONSchema, err := helpers.MultiPartFileToIFileHeader(formFile)
	if err != nil {
		return c.JSON(emsgs.OpenFileError(err).GetStatus(), emsgs.OpenFileError(err).JSON())
	}

	s := services.NewSchemaService(c)

	createdBy := c.Get(consts.ContextKeyTokenID).(string)
	vcSchemas, ierr := s.CreateFromFile(&services.SchemaCreateFromFilePayload{
		ZipJSONSchema: zipJSONSchema,
		CreatedBy:     createdBy,
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusCreated, vcSchemas)
}

func (sm SchemaController) UpdateByUpload(c core.IHTTPContext) error {
	formFile, err := c.FormFile("file")
	if err != nil {
		return c.JSON(emsgs.RequestFormError(err).GetStatus(), emsgs.RequestFormError(err).JSON())
	}

	if helpers.FileExtensionMatch(formFile.Filename, "zip") != nil {
		return c.JSON(emsgs.IsNotZipFile().GetStatus(), emsgs.IsNotZipFile().JSON())
	}
	zipJSONSchema, err := helpers.MultiPartFileToIFileHeader(formFile)
	if err != nil {
		return c.JSON(emsgs.OpenFileError(err).GetStatus(), emsgs.OpenFileError(err).JSON())
	}

	s := services.NewSchemaService(c)

	createdBy := c.Get(consts.ContextKeyTokenID).(string)
	vcSchemas, ierr := s.UpdateFromFile(c.Param("id"), &services.SchemaCreateFromFilePayload{
		ZipJSONSchema: zipJSONSchema,
		CreatedBy:     createdBy,
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, vcSchemas)
}

func (sm SchemaController) Update(c core.IHTTPContext) error {
	input := &requests.VCSchemaUpdate{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	s := services.NewSchemaService(c)

	schemaBody := json.RawMessage(utils.JSONToString(input.SchemaBody))
	createdBy := c.Get(consts.ContextKeyTokenID).(string)
	vcSchema, ierr := s.Update(c.Param("id"), &services.SchemaUpdatePayload{
		SchemaBody: &schemaBody,
		Version:    utils.GetString(input.Version),
		CreatedBy:  createdBy,
	})

	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, vcSchema)
}

func (sm SchemaController) Pagination(c core.IHTTPContext) error {

	s := services.NewSchemaService(c)
	vcSchemas, pageResponse, ierr := s.Pagination(c.QueryParam("type"), c.GetPageOptionsWithOptions(&core.PageOptionsOptions{}))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, core.NewPagination(vcSchemas, pageResponse))
}

func (sm SchemaController) Find(c core.IHTTPContext) error {
	s := services.NewSchemaService(c)
	vcSchema, ierr := s.Find(c.Param("id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, vcSchema)
}

func (sm SchemaController) FindByVersion(c core.IHTTPContext) error {
	s := services.NewSchemaService(c)
	vcSchema, ierr := s.FindByVersion(c.Param("id"), c.Param("version"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, vcSchema)
}

func (sm SchemaController) FindHistory(c core.IHTTPContext) error {
	s := services.NewSchemaService(c)
	vcSchemaHistory, ierr := s.FindHistory(c.Param("id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, vcSchemaHistory)
}

func (sm SchemaController) FindSchemaByFilename(c core.IHTTPContext) error {
	s := services.NewSchemaService(c)
	reference, ierr := s.FindReference(c.Param("id"), c.Param("version"), c.Param("schemaFilename"))
	if ierr != nil {
		if !errmsgs.IsNotFoundError(ierr) {
			return c.JSON(ierr.GetStatus(), ierr.JSON())
		}

		// This means references not found
		vcSchema, ierr := s.FindByVersion(c.Param("id"), c.Param("version"))
		if ierr != nil {
			return c.JSON(ierr.GetStatus(), ierr.JSON())
		}

		var schemaID string
		schemaIDv4 := helpers.GetJSONValue(vcSchema.SchemaBody, "id")
		if schemaIDv4 != nil {
			schemaID = schemaIDv4.(string)
		} else {
			schemaIDv6 := helpers.GetJSONValue(vcSchema.SchemaBody, "$id")
			schemaID = schemaIDv6.(string)
		}

		split := strings.Split(schemaID, "/")
		expectedSchemaFileName := split[len(split)-1]

		if c.Param("schemaFilename") != expectedSchemaFileName {
			ierr := emsgs.SchemaNotFound(vcSchema.ID, vcSchema.Version, c.Param("schemaFilename"))
			return c.JSON(ierr.GetStatus(), ierr.JSON())
		}

		return c.JSON(http.StatusOK, vcSchema.SchemaBody)
	}

	return c.JSON(http.StatusOK, reference)
}

func (sm *SchemaController) Validate(c core.IHTTPContext) error {
	input := &requests.VCSchemaValidate{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	document := json.RawMessage(utils.JSONToString(input.Document))
	s := services.NewSchemaService(c)
	valid, fields, ierr := s.ValidateDocument(utils.GetString(input.SchemaID), &document)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, views.NewVCSchemaValidateView(valid, fields))
}

func (sm *SchemaController) Types(c core.IHTTPContext) error {

	s := services.NewSchemaService(c)
	types, ierr := s.Types(c.QueryParam("type"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, core.Map{"types": types})
}
