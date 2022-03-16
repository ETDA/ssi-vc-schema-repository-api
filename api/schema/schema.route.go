package schema

import (
	"github.com/labstack/echo/v4"
	"gitlab.finema.co/finema/etda/vc-schema-api/middlewares"
	core "ssi-gitlab.teda.th/ssi/core"
)

func NewSchemaHandler(r *echo.Echo) {

	schema := &SchemaController{}

	r.POST("/schemas", core.WithHTTPContext(schema.Create), middlewares.IsAuth, middlewares.IsVCSchemaWriter)
	r.POST("/schemas/upload", core.WithHTTPContext(schema.Upload), middlewares.IsAuth, middlewares.IsVCSchemaWriter)
	r.PUT("/schemas/:id", core.WithHTTPContext(schema.Update), middlewares.IsAuth, middlewares.IsVCSchemaWriter)
	r.POST("/schemas/:id/upload", core.WithHTTPContext(schema.UpdateByUpload), middlewares.IsAuth, middlewares.IsVCSchemaWriter)
	r.GET("/schemas", core.WithHTTPContext(schema.Pagination), middlewares.IsAuth)
	r.GET("/schemas/types", core.WithHTTPContext(schema.Types), middlewares.IsAuth)
	r.GET("/schemas/:id", core.WithHTTPContext(schema.Find), middlewares.IsAuth)
	r.GET("/schemas/:id/history", core.WithHTTPContext(schema.FindHistory), middlewares.IsAuth)
	r.GET("/schemas/:id/:version", core.WithHTTPContext(schema.FindByVersion), middlewares.IsAuth)
	r.GET("/schemas/:id/:version/:schemaFilename", core.WithHTTPContext(schema.FindSchemaByFilename))
	r.POST("/schemas/validate", core.WithHTTPContext(schema.Validate)) // no middlewares for this endpoint
}
