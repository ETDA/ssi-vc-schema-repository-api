import * as Knex from "knex";
import {randomUUID, createHash} from "crypto";

export async function seed(knex: Knex): Promise<void> {
  // Check if row is empty
  const row = await knex.table('tokens').first()

  if (!row) {
    const id = randomUUID()
    const created_at = new Date()
    const role = "ADMIN"
    const token = createHash("SHA256").update(id + created_at.toString(), "utf-8").digest("hex").toString()
    console.log(`ADMIN_TOKEN: ${token}`)
    await knex("tokens").insert({
      id,
      role,
      name: "ADMIN_TOKEN",
      token,
      created_at,
      deleted_at: null
    });
  }
}
