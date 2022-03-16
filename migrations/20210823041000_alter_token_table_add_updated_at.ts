import * as Knex from "knex"

export async function up(knex: Knex): Promise<void> {
  return knex.schema.table("tokens", function (table) {
    table.dateTime("updated_at")
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable("tokens", function (table) {
    table.dropColumn("updated_at")
  })
}
