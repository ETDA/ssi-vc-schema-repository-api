package requests

import (
	"strings"

	"gitlab.finema.co/finema/etda/vc-schema-api/consts"
	core "ssi-gitlab.teda.th/ssi/core"
)

type TokenCreate struct {
	core.BaseValidator
	Name *string `json:"name"`
	Role *string `json:"role"`
}

func (r *TokenCreate) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Name, "name"))
	if r.Must(r.IsRequired(r.Role, "role")) {
		r.Must(r.IsStrIn(r.Role, strings.Join([]string{consts.TokenAdminRole, consts.TokenReadWriteRole, consts.TokenReadRole}, "|"), "role"))
	}
	return r.Error()
}
