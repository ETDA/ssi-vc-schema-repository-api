import Ajv from 'ajv'

export const schemaValidateView = (errors: Ajv['errors']) => {
  let valid: boolean
  let fields: {
    [key: string]: {
      code: string,
      message: string
    }[]
  }
  let rawError: Ajv['errors']

  valid = !errors
  if (!valid) {
    rawError = errors
    fields = {}
    rawError?.forEach((error) => {
      let key = error.instancePath.substring(1).replace(/\//g, '.')
      if (key === '') {
        key = 'schema'
      } else {
        key = 'schema.' + key
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
