package emsgs

import (
	core "ssi-gitlab.teda.th/ssi/core"
	"net/http"
)

var AuthTokenRequired = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "TOKEN_REQUIRED",
	Message: "Authorization header field is required",
}

var AuthTokenInvalid = core.Error{
	Status:  http.StatusUnauthorized,
	Code:    "INVALID_TOKEN",
	Message: "Token is invalid",
}

var PermissionDenied = core.Error{
	Status:  http.StatusForbidden,
	Code:    "PERMISSION_DENIED",
	Message: "Token does have permission to to perform this action"}
