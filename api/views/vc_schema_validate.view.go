package views

import core "ssi-gitlab.teda.th/ssi/core"

type VCSchemaValidate struct {
	Valid  bool     `json:"valid"`
	Fields core.Map `json:"fields,omitempty"`
}

func NewVCSchemaValidateView(valid bool, fields core.Map) *VCSchemaValidate {
	return &VCSchemaValidate{
		Valid:  valid,
		Fields: fields,
	}
}
