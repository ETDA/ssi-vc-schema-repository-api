export const isJoiValidationError = (e: any) => {
  return e.name === 'ValidationError' && e.isJoi
}

const isEmpty = (object: any) => {
  return Object.keys(object).length === 0
}

export const joiErrorToBadRequest = (validator: any) => {
  const errors: any = {}
  const result = validator
  if (result && result.error && result.error.details) {
    result.error.details.forEach((v: any) => {
      errors[`${v.path[0]}`] = {
        code: `${v.path.join('.')}.${v.type}${v.context.limit ? `.${v.context.limit}` : ''}`,
        message: v.message.replace(/"/g,"")
      }
    })
  }
  if (!isEmpty(errors)) {
    return {
      status: 'INVALID_PARAMS',
      message: 'Invalid parameters',
      fields: errors
    }
  }

  return null
}

export const internalServerError = (error: any) => {
  const status = "INTERNAL_SERVER_ERROR"
  const message = error.message
  if (error.stack) {
    return {
      status,
      message,
      stack: error.stack
    }
  }

  return {
    status,
    message
  }
}
