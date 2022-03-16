import Ajv2019 from "ajv/dist/2019"
import Ajv2020 from "ajv/dist/2020"
import Ajv04 from "ajv-draft-04"
import draft6MetaSchema from "ajv/dist/refs/json-schema-draft-06.json"
import draft7MetaSchema from "ajv/dist/refs/json-schema-draft-07.json"
import Ajv, { Options, ValidateFunction } from "ajv"
import { IContext } from "../types/context.type"
import axios from "axios"

interface ISchemaService {
  validateSchema(schema: any): Ajv["errors"]

  validateDocument(
    schema: any,
    document: any
  ): Promise<ValidateFunction["errors"]>
}

class schemaService implements ISchemaService {
  ctx: IContext
  ajv: Ajv2019
  ajv2020: Ajv2020
  ajv04: Ajv04

  constructor(ctx: IContext) {
    const options: Options = {
      verbose: true,
      loadSchema: this.loadSchema,
    }
    this.ctx = ctx
    this.ajv04 = new Ajv04(options)
    this.ajv04
      .addKeyword("alias")
      .addKeyword("is_date")
      .addKeyword("order")
      .addKeyword("is_group")

    this.ajv = new Ajv2019(options)
    this.ajv.addMetaSchema(draft6MetaSchema)
    this.ajv.addMetaSchema(draft7MetaSchema)
    this.ajv
      .addKeyword("alias")
      .addKeyword("is_date")
      .addKeyword("order")
      .addKeyword("is_group")

    this.ajv2020 = new Ajv2020(options)
    this.ajv2020
      .addKeyword("alias")
      .addKeyword("is_date")
      .addKeyword("order")
      .addKeyword("is_group")
  }

  validateSchema = (schema: any): Ajv["errors"] => {
    const meta: string = schema.$schema

    if (meta.includes("2020")) {
      this.ajv2020.validateSchema(schema)
      return this.ajv2020.errors
    }

    if (meta.includes("draft-04")) {
      this.ajv04.validateSchema(schema)
      return this.ajv04.errors
    }

    this.ajv.validateSchema(schema)
    return this.ajv.errors
  }

  validateDocument = async (
    schema: any,
    document: any
  ): Promise<ValidateFunction["errors"]> => {
    const meta: string = schema.$schema

    if (meta.includes("2020")) {
      const validator = await this.ajv2020.compileAsync(schema)
      validator(document)
      return validator.errors
    }

    if (meta.includes("draft-04")) {
      const validator = await this.ajv04.compileAsync(schema)
      validator(document)
      return validator.errors
    }

    const validator = await this.ajv.compileAsync(schema)
    validator(document)
    return validator.errors
  }

  loadSchema = async (uri: string) => {
    try {
      const res = await axios.get(uri)
      return res.data
    } catch (error) {
      throw error
    }
  }
}

export const newSchemaService = (ctx: IContext): ISchemaService => {
  return new schemaService(ctx)
}
