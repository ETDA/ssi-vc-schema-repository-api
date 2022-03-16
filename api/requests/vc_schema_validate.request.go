package requests

import (
	"encoding/json"

	core "ssi-gitlab.teda.th/ssi/core"
)

type VCSchemaValidate struct {
	core.BaseValidator
	SchemaID *string          `json:"schema_id"`
	Document *json.RawMessage `json:"document"`
}

func (r *VCSchemaValidate) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.SchemaID, "schema_id"))
	if r.Must(r.IsRequired(r.Document, "document")) {
		r.Must(r.IsJSONObjectNotEmpty(r.Document, "document"))
	}
	return r.Error()
}
