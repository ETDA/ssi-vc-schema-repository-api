import {ValidateFunction} from "ajv"

export const schemaValidateDocumentView = (errors: ValidateFunction['errors']) => {
  let valid: boolean
  let fields: {
    [key: string]: {
      code: string,
      message: string
    }[]
  }
  let rawError: ValidateFunction['errors']

  valid = !errors
  if (!valid) {
    rawError = errors
    fields = {}
    rawError?.forEach((error) => {
      let key = error.instancePath.substring(1).replace(/\//g, '.')
      if (key === '') {
        key = 'document'
      } else {
        key = 'document.' + key
      }

      let message = <string>error.message
      if (error.keyword == 'enum') message = `${message}: ${<string[]>(error.params.allowedValues).join(', ')}`
      if (!fields[key]) fields[key] = []
      fields[key].push(
        {
          code: error.keyword,
          message
        }
      )
    })

    return {
      valid: false,
      fields: fields,
    }
  }

  return {
    valid: true,
  }
}
