import {Express} from "express"

export interface IContext {
  app: Express
}

export class Context implements IContext {
  app: Express
  constructor(app: Express) {
    this.app = app
  }
}