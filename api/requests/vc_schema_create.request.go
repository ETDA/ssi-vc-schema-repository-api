package requests

import (
	"encoding/json"
	"strings"

	"gitlab.finema.co/finema/etda/vc-schema-api/consts"
	"gitlab.finema.co/finema/etda/vc-schema-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
)

type VCSchemaCreate struct {
	core.BaseValidator
	SchemaName *string          `json:"schema_name"`
	SchemaType *string          `json:"schema_type"`
	SchemaBody *json.RawMessage `json:"schema_body"`
}

func (r *VCSchemaCreate) Valid(ctx core.IContext) core.IError {

	r.Must(r.IsStrRequired(r.SchemaName, "schema_name"))
	r.Must(r.IsStrUnique(ctx, r.SchemaName, models.VCSchema{}.TableName(), "schema_name", "", "schema_name"))

	r.Must(r.IsStrRequired(r.SchemaType, "schema_type"))

	if r.Must(r.IsJSONStrPathRequired(r.SchemaBody, "$schema", "schema_body.$schema")) {
		r.Must(r.IsJSONPathStrIn(r.SchemaBody, "$schema", strings.Join(consts.AllowedSchemaMeta, "|"), "schema_body.$schema"))
	}

	return r.Error()
}
