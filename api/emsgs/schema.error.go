package emsgs

import (
	"fmt"
	"net/http"

	core "ssi-gitlab.teda.th/ssi/core"
)

func SchemaNotFound(id string, version string, filename string) core.Error {
	return core.Error{
		Status:  http.StatusNotFound,
		Code:    "SCHMEA_NOT_FOUND",
		Message: fmt.Sprintf("schema with id %s and version %s is does not have schema or reference with name %s", id, version, filename),
	}
}

func SchemaTypeDuplicated(schemaType string) core.Error {
	return core.Error{
		Status:  http.StatusBadRequest,
		Code:    "SCHEMA_TYPE_DUPLICATED",
		Message: fmt.Sprintf("schema type %s already exists", schemaType),
	}
}

func SchemaNameDuplicated(schemaName string) core.Error {
	return core.Error{
		Status:  http.StatusBadRequest,
		Code:    "SCHMEA_NAME_DUBPLICATED",
		Message: fmt.Sprintf("schema name %s already exists", schemaName),
	}
}

func SchemaReferenceDuplicated(id string, version string, schemaReference string) core.Error {
	return core.Error{
		Status:  http.StatusBadRequest,
		Code:    "SCHMEA_REFERENCE_DUBPLICATED",
		Message: fmt.Sprintf("schema id %s already have reference name %s for version %s", id, schemaReference, version),
	}
}

var InvalidJSONSchema = func(schemaMeta string, fields core.Map) core.Error {
	return core.Error{
		Status:  http.StatusBadRequest,
		Code:    "INVALID_JSON_SCHEMA",
		Message: fmt.Sprintf("schema is not valid according to %s", schemaMeta),
		Fields:  fields,
	}
}

var InvalidMessageInfo = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "INVALID_MESSAGE_INFO",
	Message: "messages must be array of object with following form {\"rootTagName\": string, \"versionPath\": string, \"schemaFilePath\": string, \"schematronFilePath\": string, \"vcSchema\": {\"schemaName\": string, \"schemaType\": string, \"schemaVersion\": string}\n",
}
