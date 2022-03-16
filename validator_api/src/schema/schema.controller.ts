import express from 'express'
import { IContext } from '../types/context.type'
import { schemaValidateRequest } from '../requests/schema_validate.request'
import { newSchemaService } from '../services/schema.service'
import { schemaValidateView } from '../views/schema_validate.view'
import { HTTP } from '../consts/http.const'
import {schemaValidateDocumentView } from "../views/schema_validate_decument.view"
import {internalServerError} from "../utils/error.utils"
import * as Sentry from "@sentry/node";
import {schemaValidateDocumentRequest} from "../requests/schema_document_validate.request"

class SchemaController {
  private readonly ctx: IContext

  constructor (ctx: IContext) {
    this.ctx = ctx
  }

  validateSchema: express.Handler = async (req, res): Promise<any> => {
    try {
      const error = schemaValidateRequest(req.body)
      if (error) {
        return res.status(HTTP.BAD_REQUEST).json(error)
      }
      const service = newSchemaService(this.ctx)
      const errors = await service.validateSchema(req.body.schema)
      res.json(schemaValidateView(errors))
    } catch (error) {
      Sentry.captureException(error)
      res.status(HTTP.INTERNAL_SERVER_ERROR).send(internalServerError(error))
    }
  }

  validateDocument: express.Handler = async (req, res): Promise<any> => {
    try {
      const error = schemaValidateDocumentRequest(req.body)
      if (error) {
        return res.status(HTTP.BAD_REQUEST).json(error)
      }
      const service = newSchemaService(this.ctx)
      const errors = await service.validateDocument(req.body.schema, req.body.document)
      const view = schemaValidateDocumentView(errors)
      res.json(view)
    } catch (error) {
      Sentry.captureException(error)
      res.status(HTTP.INTERNAL_SERVER_ERROR).send(internalServerError(error))
    }
  }
}

export const newSchemaController = (ctx: IContext) => new SchemaController(ctx)
