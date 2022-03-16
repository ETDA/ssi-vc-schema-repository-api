import express from "express"
import * as controllers from "./schema.controller"
import { IContext } from "../types/context.type"

export const useSchemaRouter = (
  router: express.Router,
  ctx: IContext
): void => {
  const controller = controllers.newSchemaController(ctx)
  router.post("/schema/validate-schema", controller.validateSchema)
  router.post("/schema/validate-document", controller.validateDocument)
}
