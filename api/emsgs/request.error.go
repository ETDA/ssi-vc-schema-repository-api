package emsgs

import (
	"net/http"

	core "ssi-gitlab.teda.th/ssi/core"
)

func RequestFormError(err error) core.IError {
	return core.Error{
		Status:  http.StatusBadRequest,
		Code:    "BAD_REQUEST",
		Message: err.Error(),
	}
}

func IsNotZipFile() core.IError {
	return core.Error{
		Status:  http.StatusBadRequest,
		Code:    "BAD_REQUEST",
		Message: "uploaded file must be in extension *.zip",
	}
}
