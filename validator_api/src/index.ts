import express from "express"
import * as dotenv from "dotenv"
import {useSchemaRouter} from "./schema/schema.route"
import {Context} from "./types/context.type"
import * as Sentry from "@sentry/node";



dotenv.config()
Sentry.init({dsn: process.env.SENTRY_DSN});

const PORT = process.env.PORT || "8081"
const app = express()
const router = express.Router()
app.use(express.json())

const context = new Context(app)

useSchemaRouter(router, context)
app.use(router)

const server = app.listen(PORT, () => {
  console.log(`app is listening to localhost:${PORT}`)
})


process.on("SIGTERM", async () => {
  console.log("SIGTERM detected")
  console.log("shutting down")
  server.close()
})

process.on("SIGINT", async () => {
  console.log("SIGINT detected")
  console.log("shutting down")
  server.close()
})