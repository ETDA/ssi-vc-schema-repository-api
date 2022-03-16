package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	plumbingTransportHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"gitlab.finema.co/finema/etda/vc-schema-api/consts"
	"gitlab.finema.co/finema/etda/vc-schema-api/emsgs"
	"gitlab.finema.co/finema/etda/vc-schema-api/helpers"
	"gitlab.finema.co/finema/etda/vc-schema-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"gorm.io/gorm"
)

type SchemaCreatePayload struct {
	SchemaName string
	SchemaType string
	SchemaBody *json.RawMessage
	CreatedBy  string
	Version    *string
	Filename   *string
}

type SchemaCreateFromFilePayload struct {
	ZipJSONSchema core.IFile `json:"zip_json_schema"`
	CreatedBy     string     `json:"created_by"`
}

type SchemaUpdatePayload struct {
	SchemaBody *json.RawMessage
	CreatedBy  string
	Version    string
}

type SchemaCreateResponse struct {
	Success bool `json:"success"`
}

type SchemaUpdateResponse struct {
	Success bool `json:"success"`
}

type ISchemaService interface {
	Pagination(schemaType string, options *core.PageOptions) ([]models.VCSchema, *core.PageResponse, core.IError)
	Create(payload *SchemaCreatePayload) (*models.VCSchema, core.IError)
	CreateFromFile(payload *SchemaCreateFromFilePayload) ([]models.VCSchema, core.IError)
	UpdateFromFile(schemaID string, payload *SchemaCreateFromFilePayload) (*models.VCSchema, core.IError)
	Update(schemaID string, payload *SchemaUpdatePayload) (*models.VCSchema, core.IError)
	Find(id string) (*models.VCSchema, core.IError)
	Types(schemaType string) ([]string, core.IError)
	FindByName(name string) (*models.VCSchema, core.IError)
	FindByType(vcSchemaType string) (*models.VCSchema, core.IError)
	FindByVersion(id string, version string) (*models.VCSchema, core.IError)
	FindHistory(id string) ([]models.VCSchema, core.IError)
	FindReference(id string, version string, referenceName string) (*json.RawMessage, core.IError)
	ListReferences(id string, version string) ([]json.RawMessage, core.IError)
	ValidateSchema(schema *json.RawMessage) (bool, core.Map, core.IError)
	ValidateDocument(schemaID string, document *json.RawMessage) (bool, core.Map, core.IError)
	FetchRepository()
}

type schemaService struct {
	ctx          core.IContext
	tokenService ITokenService
}

func NewSchemaService(ctx core.IContext) ISchemaService {
	return &schemaService{
		ctx:          ctx,
		tokenService: NewTokenService(ctx),
	}
}

func (s schemaService) Types(schemaType string) ([]string, core.IError) {
	var types []string
	db := s.ctx.DB()
	if schemaType != "" {
		db = core.SetSearchSimple(db, schemaType, []string{"schema_type"})
	}
	err := db.Model(&models.VCSchema{}).Pluck("schema_type", &types).Error
	if err != nil {
		return nil, s.ctx.NewError(errmsgs.DBError, errmsgs.DBError)
	}

	return helpers.UniqueStringArray(types), nil
}

func (s schemaService) Pagination(schemaType string, pageOptions *core.PageOptions) ([]models.VCSchema, *core.PageResponse, core.IError) {
	items := make([]models.VCSchema, 0)
	if len(pageOptions.OrderBy) == 0 {
		pageOptions.OrderBy = []string{"created_at DESC"}
	}

	db := s.ctx.DB()

	if schemaType != "" {
		db = db.Where("schema_type = ?", schemaType)
	}

	db = core.SetSearchSimple(db, pageOptions.Q, []string{"schema_name"})

	pageRes, err := core.Paginate(db, &items, pageOptions)
	if err != nil {
		return nil, nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return items, pageRes, nil
}

func (s schemaService) Create(payload *SchemaCreatePayload) (*models.VCSchema, core.IError) {
	valid, fields, ierr := s.ValidateSchema(payload.SchemaBody)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	if !valid {
		ierr = emsgs.InvalidJSONSchema(helpers.GetJSONValue(payload.SchemaBody, "$schema").(string), fields)
		return nil, s.ctx.NewError(ierr, ierr)
	}

	// Insert to vc_schemas table
	vcSchema, err := models.NewVCSchema(s.ctx, payload.SchemaName, payload.SchemaType, payload.Version, payload.Filename, payload.SchemaBody, payload.CreatedBy)
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.InternalServerError)
	}

	err = s.ctx.DB().Create(vcSchema).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	// Insert to vc_schema_histories table
	vcSchemaHistory := models.VCSchemaHistory{
		SchemaID:   vcSchema.ID,
		SchemaBody: vcSchema.SchemaBody,
		Version:    vcSchema.Version,
		CreatedAt:  vcSchema.CreatedAt,
		CreatedBy:  vcSchema.CreatedBy,
	}
	err = s.ctx.DB().Create(vcSchemaHistory).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return s.Find(vcSchema.ID)
}

func (s schemaService) CreateFromFile(payload *SchemaCreateFromFilePayload) ([]models.VCSchema, core.IError) {
	filename := payload.ZipJSONSchema.Name()
	err := helpers.FileExtensionMatch(filename, "zip")
	if err != nil {
		return nil, s.ctx.NewError(emsgs.OpenFileError(err), emsgs.OpenFileError(err))
	}

	readZip, err := helpers.Unzip(payload.ZipJSONSchema.Value())
	if err != nil {
		return nil, s.ctx.NewError(emsgs.OpenFileError(err), emsgs.OpenFileError(err))
	}

	rootDirectory := strings.Split(readZip.File[0].Name, "/")[0]

	readMessageInfo, err := readZip.Open(fmt.Sprintf("%s/%s", rootDirectory, "MessageInfo.json"))
	if err != nil {
		return nil, s.ctx.NewError(emsgs.OpenFileError(err), emsgs.OpenFileError(err))
	}

	rawMessageInfo, err := helpers.IOReaderToBytes(readMessageInfo)
	jsonRawMessage := json.RawMessage(rawMessageInfo)

	// by using MessageInfo.json, the message has to be in array form.
	messages, ok := helpers.GetJSONValue(&jsonRawMessage, "messages").([]interface{})
	if !ok {
		return nil, s.ctx.NewError(emsgs.InvalidMessageInfo, emsgs.InvalidMessageInfo)
	}

	createdSchemasID := make([]string, 0)
	for _, message := range messages {
		// create schema here

		// extract payload from message info
		jsonRawMessage = json.RawMessage(utils.JSONToString(message))
		schemaFilePath, ok := helpers.GetJSONValue(&jsonRawMessage, "schemaFilePath").(string)
		if !ok {
			return nil, s.ctx.NewError(emsgs.InvalidMessageInfo, emsgs.InvalidMessageInfo)
		}
		schemaName, ok := helpers.GetJSONValue(&jsonRawMessage, "vcSchema.schemaName").(string)
		if !ok {
			return nil, s.ctx.NewError(emsgs.InvalidMessageInfo, emsgs.InvalidMessageInfo)
		}
		schemaType, ok := helpers.GetJSONValue(&jsonRawMessage, "vcSchema.schemaType").(string)
		if !ok {
			return nil, s.ctx.NewError(emsgs.InvalidMessageInfo, emsgs.InvalidMessageInfo)
		}
		schemaVersion, ok := helpers.GetJSONValue(&jsonRawMessage, "vcSchema.schemaVersion").(string)
		if !ok {
			return nil, s.ctx.NewError(emsgs.InvalidMessageInfo, emsgs.InvalidMessageInfo)
		}

		// Check if schema with this name already existed
		_, ierr := s.FindByName(schemaName)
		if !errmsgs.IsNotFoundError(ierr) {
			return nil, s.ctx.NewError(emsgs.SchemaNameDuplicated(schemaName), emsgs.SchemaNameDuplicated(schemaName))
		}

		// Build schema filepath and read it
		absoluteSchemaFilePath := rootDirectory + schemaFilePath
		readSchema, err := readZip.Open(absoluteSchemaFilePath)
		if err != nil {
			return nil, s.ctx.NewError(emsgs.OpenFileError(err), emsgs.OpenFileError(err))
		}
		rawSchema, err := helpers.IOReaderToBytes(readSchema)
		if err != nil {
			return nil, s.ctx.NewError(err, emsgs.OpenFileError(err))
		}
		jsonRawMessage = rawSchema // implicit type casting

		// Get schema meta version
		meta := helpers.GetJSONValue(&jsonRawMessage, "$schema").(string)

		// Validate if schema is valid according json schema standard
		valid, fields, ierr := s.ValidateSchema(&jsonRawMessage)
		if ierr != nil {
			return nil, s.ctx.NewError(ierr, ierr)
		}

		if !valid {
			ierr = emsgs.InvalidJSONSchema(meta, fields)
			return nil, s.ctx.NewError(ierr, ierr)
		}

		// Build schema main filename
		split := strings.Split(schemaFilePath, "/")
		schemaMainFilename := split[len(split)-1]

		vcSchema, ierr := s.Create(&SchemaCreatePayload{
			SchemaName: schemaName,
			SchemaType: schemaType,
			SchemaBody: &jsonRawMessage,
			CreatedBy:  payload.CreatedBy,
			Version:    &schemaVersion,
			Filename:   &schemaMainFilename,
		})

		if ierr != nil {
			return nil, s.ctx.NewError(ierr, ierr)
		}

		// Build schemas directory to find references
		split = strings.Split(absoluteSchemaFilePath, "/")
		absoluteSchemaFileDirectory := strings.TrimSuffix(absoluteSchemaFilePath, split[len(split)-1])

		// Create reference
		for _, v := range readZip.File {

			// Filter out __MACOSX, Directory
			if !strings.Contains(v.Name, "__MACOSX") && !v.FileInfo().IsDir() && strings.Contains(v.Name, absoluteSchemaFileDirectory) && !strings.Contains(v.Name, schemaMainFilename) {
				// Build reference filename
				split := strings.Split(v.Name, "/")
				referenceFilename := split[len(split)-1]

				_, ierr := s.FindReference(vcSchema.ID, vcSchema.Version, referenceFilename)
				if !errmsgs.IsNotFoundError(ierr) {
					ierr = emsgs.SchemaReferenceDuplicated(vcSchema.ID, vcSchema.Version, referenceFilename)
					return nil, s.ctx.NewError(ierr, ierr)
				}
				// Read reference file
				readReference, err := readZip.Open(v.Name)
				if err != nil {
					return nil, s.ctx.NewError(emsgs.OpenFileError(err), emsgs.OpenFileError(err))
				}
				rawReference, err := helpers.IOReaderToBytes(readReference)
				if err != nil {
					return nil, s.ctx.NewError(err, emsgs.OpenFileError(err))
				}
				jsonRawMessage = json.RawMessage(rawReference) // implicit type casting

				// Validate if schema is valid according json schema standard
				valid, fields, ierr = s.ValidateSchema(&jsonRawMessage)
				if ierr != nil {
					return nil, s.ctx.NewError(ierr, ierr)
				}
				if !valid {
					ierr = emsgs.InvalidJSONSchema(helpers.GetJSONValue(&jsonRawMessage, "$schema").(string), fields)
					return nil, s.ctx.NewError(ierr, ierr)
				}

				newReferenceBody, ierr := s.setSchemaReferenceID(vcSchema, referenceFilename, &jsonRawMessage)
				if ierr != nil {
					return nil, s.ctx.NewError(ierr, ierr)
				}

				schemaReference := &models.VCSchemaReference{
					ID:            utils.GetUUID(),
					SchemaID:      vcSchema.ID,
					ReferenceBody: newReferenceBody,
					Name:          referenceFilename,
					CreatedAt:     utils.GetCurrentDateTime(),
					Version:       vcSchema.Version,
				}

				err = s.ctx.DB().Create(schemaReference).Error
				if err != nil {
					return nil, s.ctx.NewError(err, errmsgs.DBError)
				}
			}
		}

		createdSchemasID = append(createdSchemasID, vcSchema.ID)
	}

	createdSchemas := make([]models.VCSchema, 0)
	for _, id := range createdSchemasID {
		schema, ierr := s.Find(id)
		if err != nil {
			return nil, s.ctx.NewError(ierr, ierr)
		}
		if schema == nil {
			return nil, s.ctx.NewError(errmsgs.InternalServerError, errmsgs.InternalServerError)
		}
		createdSchemas = append(createdSchemas, *schema)
	}

	return createdSchemas, nil
}

func (s schemaService) UpdateFromFile(schemaID string, payload *SchemaCreateFromFilePayload) (*models.VCSchema, core.IError) {
	vcSchema, ierr := s.Find(schemaID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	readZip, err := helpers.Unzip(payload.ZipJSONSchema.Value())
	if err != nil {
		return nil, s.ctx.NewError(emsgs.OpenFileError(err), emsgs.OpenFileError(err))
	}

	rootDirectory := strings.Split(readZip.File[0].Name, "/")[0]
	readMessageInfo, err := readZip.Open(fmt.Sprintf("%s/%s", rootDirectory, "MessageInfo.json"))
	if err != nil {
		return nil, s.ctx.NewError(emsgs.OpenFileError(err), emsgs.OpenFileError(err))
	}

	rawMessageInfo, err := helpers.IOReaderToBytes(readMessageInfo)
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.InternalServerError)

	}
	jsonRawMessage := json.RawMessage(rawMessageInfo)

	// by using MessageInfo.json, the message has to be in array form.
	messages, ok := helpers.GetJSONValue(&jsonRawMessage, "messages").([]interface{})
	if !ok {
		return nil, s.ctx.NewError(emsgs.InvalidMessageInfo, emsgs.InvalidMessageInfo)
	}

	// Bulk update is not supported
	if len(messages) != 1 {
		return nil, s.ctx.NewError(emsgs.InvalidMessageInfo, emsgs.InvalidMessageInfo)
	}

	message := messages[0]
	// extract payload from message info
	jsonRawMessage = json.RawMessage(utils.JSONToString(message))
	schemaFilePath, ok := helpers.GetJSONValue(&jsonRawMessage, "schemaFilePath").(string)
	if !ok {
		return nil, s.ctx.NewError(emsgs.InvalidMessageInfo, emsgs.InvalidMessageInfo)
	}
	schemaVersion, ok := helpers.GetJSONValue(&jsonRawMessage, "vcSchema.schemaVersion").(string)
	if !ok {
		return nil, s.ctx.NewError(emsgs.InvalidMessageInfo, emsgs.InvalidMessageInfo)
	}

	_, valiatorError := helpers.IsValidVersionUpdate(&schemaVersion, vcSchema.Version, "messages[0].vcSchema.schemaVersion")
	if valiatorError != nil {
		ierr = core.Error{
			Status:  http.StatusBadRequest,
			Code:    valiatorError.Code,
			Message: valiatorError.Message,
			Data:    valiatorError.Data,
		}
		return nil, s.ctx.NewError(ierr, ierr)
	}

	// Build schema filepath and read it
	absoluteSchemaFilePath := rootDirectory + schemaFilePath
	readSchema, err := readZip.Open(absoluteSchemaFilePath)
	if err != nil {
		return nil, s.ctx.NewError(emsgs.OpenFileError(err), emsgs.OpenFileError(err))
	}
	rawSchema, err := helpers.IOReaderToBytes(readSchema)
	if err != nil {
		return nil, s.ctx.NewError(err, emsgs.OpenFileError(err))
	}
	jsonRawMessage = rawSchema // implicit type casting

	// Get schema meta version
	meta := helpers.GetJSONValue(&jsonRawMessage, "$schema").(string)

	// Validate if schema is valid according json schema standard
	valid, fields, ierr := s.ValidateSchema(&jsonRawMessage)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	if !valid {
		ierr = emsgs.InvalidJSONSchema(meta, fields)
		return nil, s.ctx.NewError(ierr, ierr)
	}

	vcSchema, ierr = s.Update(vcSchema.ID, &SchemaUpdatePayload{
		SchemaBody: &jsonRawMessage,
		Version:    schemaVersion,
		CreatedBy:  payload.CreatedBy,
	})

	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	// Build schema main filename
	split := strings.Split(schemaFilePath, "/")
	schemaMainFilename := split[len(split)-1]

	// Build schemas directory to find references
	split = strings.Split(absoluteSchemaFilePath, "/")
	absoluteSchemaFileDirectory := strings.TrimSuffix(absoluteSchemaFilePath, split[len(split)-1])

	// Create reference
	for _, v := range readZip.File {

		// Filter out __MACOSX, Directory
		if !strings.Contains(v.Name, "__MACOSX") && !v.FileInfo().IsDir() && strings.Contains(v.Name, absoluteSchemaFileDirectory) && !strings.Contains(v.Name, schemaMainFilename) {
			// Build reference filename
			split := strings.Split(v.Name, "/")
			referenceFilename := split[len(split)-1]

			_, ierr := s.FindReference(vcSchema.ID, vcSchema.Version, referenceFilename)
			if !errmsgs.IsNotFoundError(ierr) {
				ierr = emsgs.SchemaReferenceDuplicated(vcSchema.ID, vcSchema.Version, referenceFilename)
				return nil, s.ctx.NewError(ierr, ierr)
			}
			// Read reference file
			readReference, err := readZip.Open(v.Name)
			if err != nil {
				return nil, s.ctx.NewError(emsgs.OpenFileError(err), emsgs.OpenFileError(err))
			}
			rawReference, err := helpers.IOReaderToBytes(readReference)
			if err != nil {
				return nil, s.ctx.NewError(err, emsgs.OpenFileError(err))
			}
			jsonRawMessage = json.RawMessage(rawReference) // implicit type casting

			// Validate if schema is valid according json schema standard
			valid, fields, ierr = s.ValidateSchema(&jsonRawMessage)
			if ierr != nil {
				return nil, s.ctx.NewError(ierr, ierr)
			}
			if !valid {
				ierr = emsgs.InvalidJSONSchema(helpers.GetJSONValue(&jsonRawMessage, "$schema").(string), fields)
				return nil, s.ctx.NewError(ierr, ierr)
			}

			newReferenceBody, ierr := s.setSchemaReferenceID(vcSchema, referenceFilename, &jsonRawMessage)
			if ierr != nil {
				return nil, s.ctx.NewError(ierr, ierr)
			}

			schemaReference := &models.VCSchemaReference{
				ID:            utils.GetUUID(),
				SchemaID:      vcSchema.ID,
				ReferenceBody: newReferenceBody,
				Name:          referenceFilename,
				CreatedAt:     utils.GetCurrentDateTime(),
				Version:       vcSchema.Version,
			}

			err = s.ctx.DB().Create(schemaReference).Error
			if err != nil {
				return nil, s.ctx.NewError(err, errmsgs.DBError)
			}
		}
	}

	return s.Find(vcSchema.ID)
}

func (s schemaService) Update(schemaID string, payload *SchemaUpdatePayload) (*models.VCSchema, core.IError) {
	valid, fields, ierr := s.ValidateSchema(payload.SchemaBody)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	if !valid {
		ierr = emsgs.InvalidJSONSchema(helpers.GetJSONValue(payload.SchemaBody, "$schema").(string), fields)
		return nil, s.ctx.NewError(ierr, ierr)
	}

	vcSchema, ierr := s.Find(schemaID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	vcSchema.Version = payload.Version
	vcSchema.CreatedBy = payload.CreatedBy
	vcSchema.UpdatedAt = utils.GetCurrentDateTime()

	newSchemaBody, ierr := s.setSchemaReferenceID(vcSchema, fmt.Sprintf("%s.json", helpers.AlphaNumericize(vcSchema.SchemaName)), payload.SchemaBody)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	vcSchema.SchemaBody = newSchemaBody

	// Update to vc_schema_table
	err := s.ctx.DB().Updates(vcSchema).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	// Insert to vc_schema_histories table
	vcSchemaHistory := models.VCSchemaHistory{
		SchemaID:   vcSchema.ID,
		SchemaBody: newSchemaBody,
		Version:    vcSchema.Version,
		CreatedAt:  vcSchema.UpdatedAt,
		CreatedBy:  vcSchema.CreatedBy,
	}
	err = s.ctx.DB().Create(vcSchemaHistory).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return s.Find(vcSchema.ID)
}

func (s schemaService) Find(id string) (*models.VCSchema, core.IError) {
	item := &models.VCSchema{}

	err := s.ctx.DB().Where("id = ?", id).First(item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(errmsgs.NotFound, errmsgs.NotFound)
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return item, nil
}

func (s schemaService) FindByName(name string) (*models.VCSchema, core.IError) {
	item := &models.VCSchema{}

	err := s.ctx.DB().Where("schema_name = ?", name).First(item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(errmsgs.NotFound, errmsgs.NotFound)
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return item, nil
}

func (s schemaService) FindByType(vcSchemaType string) (*models.VCSchema, core.IError) {
	item := &models.VCSchema{}

	err := s.ctx.DB().Where("schema_type = ?", vcSchemaType).First(item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(errmsgs.NotFound, errmsgs.NotFound)
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return item, nil
}

func (s schemaService) FindByVersion(id string, version string) (*models.VCSchema, core.IError) {
	item := &models.VCSchemaHistory{}

	err := s.ctx.DB().Preload("Schema").Where("schema_id = ? AND version = ?", id, version).First(item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(errmsgs.NotFound, errmsgs.NotFound)
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	schema := item.Schema
	schema.SchemaBody = item.SchemaBody
	schema.Version = item.Version
	schema.CreatedAt = item.CreatedAt

	return schema, nil
}

func (s schemaService) FindHistory(id string) ([]models.VCSchema, core.IError) {
	_, ierr := s.Find(id)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	items := make([]models.VCSchemaHistory, 0)
	err := s.ctx.DB().Preload("Schema").Where("schema_id = ?", id).Order("created_at DESC").Find(&items).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	vcSchemas := make([]models.VCSchema, 0)
	for _, item := range items {
		schema := item.Schema
		schema.SchemaBody = item.SchemaBody
		schema.Version = item.Version
		schema.CreatedAt = item.CreatedAt

		vcSchemas = append(vcSchemas, *schema)
	}

	return vcSchemas, nil
}

func (s schemaService) FindReference(id string, version string, referenceName string) (*json.RawMessage, core.IError) {
	_, ierr := s.Find(id)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	item := &models.VCSchemaReference{}
	err := s.ctx.DB().Where("schema_id = ? AND version = ? AND name = ?", id, version, referenceName).First(item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(errmsgs.NotFound, errmsgs.NotFound)
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return item.ReferenceBody, nil
}

func (s schemaService) ListReferences(id string, version string) ([]json.RawMessage, core.IError) {
	_, ierr := s.Find(id)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	items := make([]models.VCSchemaReference, 0)

	err := s.ctx.DB().Where("schema_id = ? AND version = ?", id, version).Find(&items).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	jsonMessageArray := make([]json.RawMessage, 0)

	for _, item := range items {
		jsonMessageArray = append(jsonMessageArray, *item.ReferenceBody)
	}

	return jsonMessageArray, nil
}

type validateResult struct {
	Valid  bool     `json:"valid"`
	Fields core.Map `json:"fields"`
}

func (s schemaService) ValidateSchema(schema *json.RawMessage) (bool, core.Map, core.IError) {
	utils.LogStruct(schema)
	res := &validateResult{}
	ierr := core.RequesterToStruct(res, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post(
			"/schema/validate-schema",
			core.Map{"schema": schema},
			&core.RequesterOptions{
				BaseURL: s.ctx.ENV().String(consts.ENVValidatorAPIEndpoint),
				Headers: http.Header{
					"content-type": []string{"application/json"},
				},
			},
		)
	})

	if ierr != nil {
		return false, nil, s.ctx.NewError(ierr, ierr)
	}

	return res.Valid, res.Fields, nil
}

func (s schemaService) ValidateDocument(schemaID string, document *json.RawMessage) (bool, core.Map, core.IError) {
	var schemaMapper map[string]interface{}
	ierr := core.RequesterToStruct(&schemaMapper, func() (*core.RequestResponse, error) {
		cc := s.ctx.(core.IHTTPContext)
		return s.ctx.Requester().Get(schemaID, &core.RequesterOptions{
			Headers: cc.Request().Header,
		})
	})
	if ierr != nil {
		return false, nil, s.ctx.NewError(ierr, ierr)
	}

	schema := json.RawMessage(utils.JSONToString(schemaMapper))

	res := &validateResult{}
	ierr = core.RequesterToStruct(res, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post(
			"/schema/validate-document",
			core.Map{
				"schema":   &schema,
				"document": document,
			},
			&core.RequesterOptions{
				BaseURL: s.ctx.ENV().String(consts.ENVValidatorAPIEndpoint),
				Headers: http.Header{
					"content-type": []string{"application/json"},
				},
			},
		)
	})

	if ierr != nil {
		return false, nil, s.ctx.NewError(ierr, ierr)
	}

	return res.Valid, res.Fields, nil
}

func (s schemaService) setSchemaReferenceID(vcSchema *models.VCSchema, filename string, jsonRawMessage *json.RawMessage) (*json.RawMessage, core.IError) {
	var newReferenceBody *json.RawMessage
	var err error
	if meta, ok := helpers.GetJSONValue(jsonRawMessage, "$schema").(string); ok && strings.Contains(meta, "draft-04") {
		newSchemaID := models.NewVCSchemaID(s.ctx, vcSchema.ID, vcSchema.Version, filename)
		newReferenceBody, err = helpers.SetJSONValue(jsonRawMessage, "id", newSchemaID)
		if err != nil {
			return nil, s.ctx.NewError(err, errmsgs.InternalServerError)
		}
	} else {
		newSchemaID := models.NewVCSchemaID(s.ctx, vcSchema.ID, vcSchema.Version, filename)
		newReferenceBody, err = helpers.SetJSONValue(jsonRawMessage, "$id", newSchemaID)
		if err != nil {
			return nil, s.ctx.NewError(err, errmsgs.InternalServerError)
		}
	}

	return newReferenceBody, nil
}

type messageInfo struct {
	schemaName    string
	schemaType    string
	schemaVersion string
}

func (s schemaService) readMessageInfo(file core.IFile) (*messageInfo, core.IError) {
	filename := file.Name()
	err := helpers.FileExtensionMatch(filename, "zip")
	if err != nil {
		return nil, s.ctx.NewError(emsgs.OpenFileError(err), emsgs.OpenFileError(err))
	}

	readZip, err := helpers.Unzip(file.Value())
	if err != nil {
		return nil, s.ctx.NewError(emsgs.OpenFileError(err), emsgs.OpenFileError(err))
	}

	rootDirectory := strings.Split(readZip.File[0].Name, "/")[0]

	readMessageInfo, err := readZip.Open(fmt.Sprintf("%s/%s", rootDirectory, "MessageInfo.json"))
	if err != nil {
		return nil, s.ctx.NewError(emsgs.OpenFileError(err), emsgs.OpenFileError(err))
	}

	rawMessageInfo, err := helpers.IOReaderToBytes(readMessageInfo)
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.InternalServerError)

	}
	jsonRawMessage := json.RawMessage(rawMessageInfo)

	// by using MessageInfo.json, the message has to be in array form.
	messages, ok := helpers.GetJSONValue(&jsonRawMessage, "messages").([]interface{})
	if !ok {
		return nil, s.ctx.NewError(emsgs.InvalidMessageInfo, emsgs.InvalidMessageInfo)
	}

	// not support bulk reading at this time
	message := messages[0]
	jsonRawMessage = json.RawMessage(utils.JSONToString(message))
	if !ok {
		return nil, s.ctx.NewError(emsgs.InvalidMessageInfo, emsgs.InvalidMessageInfo)
	}
	schemaName, ok := helpers.GetJSONValue(&jsonRawMessage, "vcSchema.schemaName").(string)
	if !ok {
		return nil, s.ctx.NewError(emsgs.InvalidMessageInfo, emsgs.InvalidMessageInfo)
	}
	schemaType, ok := helpers.GetJSONValue(&jsonRawMessage, "vcSchema.schemaType").(string)
	if !ok {
		return nil, s.ctx.NewError(emsgs.InvalidMessageInfo, emsgs.InvalidMessageInfo)
	}
	schemaVersion, ok := helpers.GetJSONValue(&jsonRawMessage, "vcSchema.schemaVersion").(string)
	if !ok {
		return nil, s.ctx.NewError(emsgs.InvalidMessageInfo, emsgs.InvalidMessageInfo)
	}

	return &messageInfo{
		schemaName:    schemaName,
		schemaType:    schemaType,
		schemaVersion: schemaVersion,
	}, nil
}

func (s schemaService) FetchRepository() {
	fetch := func() {
		path := "/tmp/vcschema/medical-certificate"
		_, err := git.PlainClone(path, false, &git.CloneOptions{
			URL:      "https://gitlab.com/quazzer/test-schema.git",
			Progress: os.Stdout,
			Auth: &plumbingTransportHttp.BasicAuth{
				Username: s.ctx.ENV().String(consts.ENVGitRepositoryUsername),
				Password: s.ctx.ENV().String(consts.ENVGitRepositoryPassword),
			},
		})
		if err != nil && err != git.ErrRepositoryAlreadyExists {
			s.ctx.Log().Error(err)
			os.Exit(1)
		}

		if err == git.ErrRepositoryAlreadyExists {
			repo, _err := git.PlainOpen(path)
			err = _err
			if err != nil {
				s.ctx.Log().Error(err)
				os.Exit(1)
			}

			workTree, _err := repo.Worktree()
			err = _err
			if err != nil {
				s.ctx.Log().Error(err)
				os.Exit(1)
			}

			err = workTree.Pull(&git.PullOptions{
				Progress: os.Stdout,
				Auth: &plumbingTransportHttp.BasicAuth{
					Username: s.ctx.ENV().String(consts.ENVGitRepositoryUsername),
					Password: s.ctx.ENV().String(consts.ENVGitRepositoryPassword),
				},
			})
			if err != nil && err != git.NoErrAlreadyUpToDate {
				s.ctx.Log().Error(err)
				os.Exit(1)
			}
		}

		if err == git.NoErrAlreadyUpToDate {
			log.Println("repository already up-to-date")
			return
		}

		token, ierr := s.tokenService.FindAdminToken()
		if ierr != nil {
			s.ctx.NewError(ierr, ierr)
			os.Exit(1)
		}

		zipFile, err := helpers.ZipFromDirectory(path)
		if err != nil {
			s.ctx.Log().Error(err)
			os.Exit(1)
		}

		file := core.NewFile("MedicalCertificate.zip", zipFile)

		info, ierr := s.readMessageInfo(file)
		if ierr != nil {
			s.ctx.NewError(ierr, ierr)
			os.Exit(1)
		}

		schema, ierr := s.FindByName(info.schemaName)
		if ierr != nil && !errmsgs.IsNotFoundError(ierr) {
			s.ctx.NewError(ierr, ierr)
			os.Exit(1)
		}

		if schema != nil {
			// means schema already exist, call update service
			_, ierr := s.UpdateFromFile(schema.ID, &SchemaCreateFromFilePayload{
				CreatedBy:     token.ID,
				ZipJSONSchema: file,
			})
			if ierr != nil {
				s.ctx.NewError(ierr, ierr)
				os.Exit(1)
			}
			return
		}

		_, ierr = s.CreateFromFile(&SchemaCreateFromFilePayload{
			CreatedBy:     token.ID,
			ZipJSONSchema: file,
		})
		if ierr != nil {
			s.ctx.NewError(ierr, ierr)
			os.Exit(1)
		}
	}
	now := time.Now()
	// 02:00 GMT+7 = 19:00 UTC
	schedule := time.Date(now.Year(), now.Month(), now.Day(), 8, 35, 0, 0, time.UTC)

	helpers.RunAtTime(fetch, &schedule, 24*time.Hour, false)
}
