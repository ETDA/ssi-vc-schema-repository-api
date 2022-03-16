package requests

import (
	"encoding/json"
	"strings"

	"gitlab.finema.co/finema/etda/vc-schema-api/consts"
	"gitlab.finema.co/finema/etda/vc-schema-api/helpers"
	"gitlab.finema.co/finema/etda/vc-schema-api/services"
	core "ssi-gitlab.teda.th/ssi/core"
)

type VCSchemaUpdate struct {
	core.BaseValidator
	SchemaBody *json.RawMessage `json:"schema_body"`
	Version    *string          `json:"version"`
}

func (r *VCSchemaUpdate) Valid(ctx core.IContext) core.IError {
	c := ctx.(core.IHTTPContext)

	if r.Must(r.IsJSONStrPathRequired(r.SchemaBody, "$schema", "schema_body.$schema")) {
		r.Must(r.IsJSONPathStrIn(r.SchemaBody, "$schema", strings.Join(consts.AllowedSchemaMeta, "|"), "schema_body.$schema"))
	}
	schemaSvc := services.NewSchemaService(ctx)
	vcSchema, _ := schemaSvc.Find(c.Param("id"))
	// incorrect id give nil pointer
	// if the database is down, it already showed id does not exist from IsMongoExistsWithCondition validator.
	if vcSchema != nil {
		if r.Must(r.IsStrRequired(r.Version, "version")) {
			r.Must(helpers.IsValidVersionUpdate(r.Version, vcSchema.Version, "schema_body.version"))
		}
	}
	return r.Error()
}
