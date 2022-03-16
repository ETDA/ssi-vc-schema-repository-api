package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gitlab.finema.co/finema/etda/vc-schema-api/consts"
	"gitlab.finema.co/finema/etda/vc-schema-api/helpers"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"gorm.io/gorm"
)

type VCSchema struct {
	ID         string           `json:"id" gorm:"id"`
	SchemaName string           `json:"schema_name" gorm:"schema_name"`
	SchemaType string           `json:"schema_type" gorm:"schema_type"`
	SchemaBody *json.RawMessage `json:"schema_body" gorm:"schema_body"`
	Version    string           `json:"version" gorm:"version"`
	CreatedBy  string           `json:"created_by" gorm:"created_by"`
	CreatedAt  *time.Time       `json:"created_at" gorm:"created_at"`
	UpdatedAt  *time.Time       `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt  *gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"deleted_at"`
}

func (m VCSchema) TableName() string {
	return "vc_schemas"
}

func NewVCSchema(ctx core.IContext, schemaName string, schemaType string, schemaVersion *string, schemaFilename *string, schemaBody *json.RawMessage, createdAtBy string) (*VCSchema, error) {
	version := "1.0.0"
	filename := fmt.Sprintf("%s.json", helpers.AlphaNumericize(schemaName))
	if schemaVersion != nil {
		version = utils.GetString(schemaVersion)
	}
	if schemaFilename != nil {
		filename = utils.GetString(schemaFilename)
	}

	id := utils.GetUUID()
	schemaID := NewVCSchemaID(ctx, id, version, filename)
	var newSchemaBody *json.RawMessage
	var err error
	if meta, ok := helpers.GetJSONValue(schemaBody, "$schema").(string); ok && strings.Contains(meta, "draft-04") {
		newSchemaBody, err = helpers.SetJSONValue(schemaBody, "id", schemaID)
		if err != nil {
			return nil, err
		}
	} else {
		newSchemaBody, err = helpers.SetJSONValue(schemaBody, "$id", schemaID)
		if err != nil {
			return nil, err
		}
	}

	return &VCSchema{
		ID:         id,
		SchemaName: schemaName,
		SchemaType: schemaType,
		SchemaBody: newSchemaBody,
		Version:    version,
		CreatedBy:  createdAtBy,
		CreatedAt:  utils.GetCurrentDateTime(),
		UpdatedAt:  utils.GetCurrentDateTime(),
	}, nil
}

func NewVCSchemaID(ctx core.IContext, id string, version string, filename string) string {
	return strings.Join([]string{ctx.ENV().String(consts.ENVVCSchemaRepositoryBaseURL), id, version, filename}, "/")
}
