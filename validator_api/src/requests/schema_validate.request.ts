import Joi from 'joi'
import { allowedSchemaMeta } from '../consts/schema.const'
import { joiErrorToBadRequest } from '../utils/error.utils'

export const schemaValidateRequest = (body: any) => joiErrorToBadRequest(
  Joi.object({
    schema: Joi.object({ $schema: Joi.string().valid(...allowedSchemaMeta).required() }).required().unknown(true)
  }
  ).validate(body)
)

