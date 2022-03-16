package emsgs

import (
	"net/http"

	core "ssi-gitlab.teda.th/ssi/core"
)

func OpenFileError(err error) core.IError {
	return core.Error{
		Status:  http.StatusBadRequest,
		Code:    "FILE_CANNOT_BE_OPENED",
		Message: err.Error(),
	}
}
