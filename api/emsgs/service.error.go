package emsgs

import (
	core "ssi-gitlab.teda.th/ssi/core"
	"net/http"
)

var ExternalServiceError = core.Error{
	Status:  http.StatusBadGateway,
	Code:    "BAD_GATEWAY",
	Message: "External service error",
}
