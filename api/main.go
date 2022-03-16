package main

import (
	"fmt"
	"os"

	"gitlab.finema.co/finema/etda/vc-schema-api/home"
	"gitlab.finema.co/finema/etda/vc-schema-api/schema"
	"gitlab.finema.co/finema/etda/vc-schema-api/services"
	"gitlab.finema.co/finema/etda/vc-schema-api/token"
	core "ssi-gitlab.teda.th/ssi/core"
)

func main() {
	env := core.NewEnv()

	mysql, err := core.NewDatabase(env.Config()).Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "MySQL: %v", err)
		os.Exit(1)
	}
	contextOption := &core.ContextOptions{
		DB:  mysql,
		ENV: env,
	}
	e := core.NewHTTPServer(&core.HTTPContextOptions{
		ContextOptions: contextOption,
	})

	ctx := core.NewContext(contextOption)
	schemaService := services.NewSchemaService(ctx)
	go schemaService.FetchRepository()

	home.NewHomeHandler(e)
	schema.NewSchemaHandler(e)
	token.NewSchemaHandler(e)

	core.StartHTTPServer(e, env)
}
