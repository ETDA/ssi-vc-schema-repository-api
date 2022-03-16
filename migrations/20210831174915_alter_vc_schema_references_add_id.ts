import * as Knex from "knex"

export async function up(knex: Knex): Promise<void> {
  return knex.schema.table('vc_schema_references', function (table) {
    table.uuid('id').primary()
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('vc_schema_references', function (table) {
    table.dropColumn('id')
  })
}
